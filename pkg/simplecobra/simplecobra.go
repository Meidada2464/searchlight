/**
 * Package simplecobra
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/14 23:03
 * param this pkg is copy form https://github.com/spf13/cobra, it a simple cobra architecture
 */

package simplecobra

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type (
	Commander interface {
		// Name 初始化命令的名字
		Name() string

		// Init 这是添加标志、短描述和长描述等的地方
		Init(*Commandeer) error

		PreRun(this, runner *Commandeer) error

		Run(ctx context.Context, cd *Commandeer, args []string) error

		Commands() []Commander
	}

	// Commandeer 最终用到的结构体的包装类,包装了两层
	Commandeer struct {
		Command      Commander
		CobraCommand *cobra.Command

		Root        *Commandeer
		Parent      *Commandeer
		commandeers []*Commandeer
	}

	Exec struct {
		c *Commandeer
	}

	runErr struct {
		error
	}

	CommandError struct {
		Err error
	}
)

// New 初始化一个RootCmd，接收接口类型参数，返回一个Exec
func New(rootCmd Commander) (*Exec, error) {
	rootCd := &Commandeer{
		Command: rootCmd,
	}

	// 递归添加所有的命令,将统一级别的子命令放在同层后递归添加
	var addCommands func(cd *Commandeer, cmd Commander)
	addCommands = func(cd *Commandeer, cmd Commander) {
		cd2 := &Commandeer{
			Root:    rootCd,
			Parent:  cd,
			Command: cmd,
		}

		cd.commandeers = append(cd.commandeers, cd2)

		// 如果子命令还有子命令，则再递归
		for _, c := range cmd.Commands() {
			addCommands(cd2, c)
		}
	}

	for _, cmd := range rootCmd.Commands() {
		addCommands(rootCd, cmd)
	}

	// 校验结构体是否有错误
	if err := rootCd.compile(); err != nil {
		return nil, err
	}

	// 最终返回的是一个树状结构体
	return &Exec{
		c: rootCd,
	}, nil
}

func (c *Commandeer) compile() error {
	useCommandFlagsArgs := "[command] [flags]"

	// 如果没有子命令
	if len(c.commandeers) == 0 {
		useCommandFlagsArgs = "[flags] [args]"
	}

	// 初始化一个根命令
	// TODO 这里的useCommandFlagsArgs是什么意思？
	c.CobraCommand = &cobra.Command{
		Use: fmt.Sprintf("%s %s", c.Command.Name(), useCommandFlagsArgs),
		// 运行前错误检查
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.Command.Run(cmd.Context(), c, args); err != nil {
				return &runErr{
					error: err,
				}
			}
			return nil
		},

		// 预运行前检查
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.init()
		},
		SilenceErrors:              true,
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 2,
	}
	// 这是添加标志、短描述和长描述等的地方
	if err := c.Command.Init(c); err != nil {
		return err
	}

	// 递归添加所有的命令
	for _, cc := range c.commandeers {
		if err := cc.compile(); err != nil {
			return err
		}
		c.CobraCommand.AddCommand(cc.CobraCommand)
	}
	return nil
}

func (c *Commandeer) init() error {
	var ancestors []*Commandeer
	{
		cd := c
		for cd != nil {
			ancestors = append(ancestors, cd)
			cd = cd.Parent
		}
	}

	// 从根开始初始化所有这些命令,跑PreRun，看是否有问题
	for i := len(ancestors) - 1; i >= 0; i-- {
		cd := ancestors[i]
		if err := cd.Command.PreRun(cd, c); err != nil {
			return err
		}
	}
	return nil
}

// Execute 根据输入参数，找到最适合的命令执行
// 执行从根命令开始执行命令树。 args 通常用 os.Args[1:] 填充。
func (r *Exec) Execute(ctx context.Context, args []string) (*Commandeer, error) {
	if args == nil {
		// 如果 args 为零，Cobra 会回退到 os.Args[1:]
		args = []string{}
	}

	// 设置输入参数,带上ctx运行程序，做超时控制等
	r.c.CobraCommand.SetArgs(args)
	cobraCommand, err := r.c.CobraCommand.ExecuteContextC(ctx)
	var cd *Commandeer
	if cobraCommand != nil {
		if err == nil {
			err = checkArgs(cobraCommand, args)
		}

		//	找到可以执行的commandeer,递归查找
		var find func(command *cobra.Command, commandeer *Commandeer) *Commandeer
		find = func(what *cobra.Command, in *Commandeer) *Commandeer {
			if in.CobraCommand == what {
				return in
			}
			for _, in2 := range in.commandeers {
				if found := find(what, in2); found != nil {
					return found
				}
			}
			return nil
		}
		cd = find(cobraCommand, r.c)
	}
	return cd, wrapErr(err)
}

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if rerr, ok := err.(*runErr); ok {
		return rerr.error
	}
	return &CommandError{Err: err}
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("command error: %v", e.Err)
}

// Is reports whether e is of type *CommandError.
func (*CommandError) Is(e error) bool {
	_, ok := e.(*CommandError)
	return ok
}

// IsCommandError  reports whether any error in err's tree matches CommandError.
func IsCommandError(err error) bool {
	return errors.Is(err, &CommandError{})
}

func checkArgs(cmd *cobra.Command, args []string) error {
	// 没有子命令就返回nil，有则检查args
	if !cmd.HasSubCommands() {
		return nil
	}

	var commandName string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			break
		}
		commandName = arg
	}

	if commandName == "" || cmd.Name() == commandName {
		return nil
	}

	// 检查是否有别名
	if cmd.HasAlias(commandName) {
		return nil
	}
	return fmt.Errorf("unknow command %q for %q%s", args[1], cmd.CommandPath(), findSuggestions(cmd, commandName))
}

func findSuggestions(cmd *cobra.Command, arg string) string {
	if cmd.DisableSuggestions {
		return ""
	}
	suggestionsString := ""
	if suggestions := cmd.SuggestionsFor(arg); len(suggestions) > 0 {
		suggestionsString = "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		}
	}
	return suggestionsString
}
