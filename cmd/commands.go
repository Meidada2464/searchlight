/**
 * Package cmd
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/11 23:43
 */

package cmd

import (
	"context"
	"errors"
	"fmt"
	"searchlight/pkg/simplecobra"
)

var errHelp = errors.New("help requested")

// Register other subcommands
func newExec() (*simplecobra.Exec, error) {
	rootCmd := &rootCommand{
		commands: []simplecobra.Commander{},
	}

	return simplecobra.New(rootCmd)
}

func Execute(args []string) error {
	exec, err := newExec()
	if err != nil {
		fmt.Println("commands pkg Execute error: ", err)
		return err
	}

	cd, err := exec.Execute(context.Background(), args)
	if err != nil {
		return err
	}

	if err != nil {
		if err == errHelp {
			cd.CobraCommand.Help()
			fmt.Println()
			return nil
		}
		if simplecobra.IsCommandError(err) {
			// Print the help, but also return the error to fail the command.
			cd.CobraCommand.Help()
			fmt.Println()
		}
	}
	return err
}
