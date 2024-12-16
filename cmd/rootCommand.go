/**
 * Package cmd
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/16 11:19
 */

package cmd

import (
	"context"
	"fmt"
	"searchlight/pkg/simplecobra"
)

type (
	rootCommand struct {
		commands []simplecobra.Commander

		// Flags
		source string
	}
)

func (r *rootCommand) Name() string {
	return "light"
}

func (r *rootCommand) Init(cd *simplecobra.Commandeer) error {
	return r.initRootCommand("", cd)
}

func (r *rootCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (r *rootCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return nil
}

func (r *rootCommand) Commands() []simplecobra.Commander {
	return nil
}

// 初始化子任务等命令
func (r *rootCommand) initRootCommand(subCommandName string, cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	commandName := "light"
	if subCommandName != "" {
		commandName = subCommandName
	}

	cmd.Use = fmt.Sprintf("%s [flags]", commandName)
	cmd.Short = "The cli tool is used to solve daily network problems!"
	cmd.Long = "searchlight is a cli tool is used to solve daily network problems!"

	// in their can config persistent flags and local flags
	cmd.PersistentFlags().StringVarP(&r.source, "source", "s", "", "The source of searchlight")
	_ = cmd.MarkFlagDirname("source")
	return nil
}
