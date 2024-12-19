/**
 * Package cmd
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/18 10:47
 */

package cmd

import (
	"context"
	"errors"
	"searchlight/internal/cmdPing"
	"searchlight/pkg/nic"
	"searchlight/pkg/simplecobra"
)

type (
	cmdPingCommand struct {
		commands []simplecobra.Commander

		// localFlags
		// source ip
		srcIp string
		// target ip
		tarIp string
		// Set the number of times a request response is completed
		count int
		// Set the packet size
		size int
		// Set the interval time
		interVal int
		// Set the timeout
		timeOut int
		// Set the ping pattern
		pattern bool
	}
)

func newCmdPingCommand() *cmdPingCommand {
	c := &cmdPingCommand{}
	return c
}

func (c *cmdPingCommand) Name() string {
	return "cmdPing"
}

func (c *cmdPingCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "ping detection using local ping tool"
	cmd.Long = "cmdPing is equivalent to directly invoking the local ping tool for network detection of the target. Note that the ping tool is installed on the machine"
	cmd.Aliases = []string{"cp", "cping", "cPing", "cmdPing"}

	cmd.Flags().StringVarP(&c.srcIp, "srcIP", "I", "", "The detect of source ip or local ip")
	cmd.Flags().StringVarP(&c.tarIp, "targetIP", "T", "", "The detect of target ip or remote ip")
	cmd.Flags().IntVarP(&c.count, "count", "c", 8, "Set the number of times a request response is completed,default 8")
	cmd.Flags().IntVarP(&c.size, "size", "s", 24, "Set the packet size,default 24 bytes")
	cmd.Flags().IntVarP(&c.interVal, "interval", "i", 1, "Set the interval time,default 1 second")
	cmd.Flags().IntVarP(&c.timeOut, "timeout", "t", 30, "Set the timeout,default 30 second")
	cmd.Flags().BoolVarP(&c.pattern, "pattern", "p", false, "Set the ping pattern, detailed data is displayed by default")
	return nil
}

func (c *cmdPingCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (c *cmdPingCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	// Check data
	if c.tarIp == "" {
		return errors.New("target ip is empty, please set it. e.g. -T 8.8.8.8")
	}

	if nic.GetIpType(c.srcIp) != nic.GetIpType(c.tarIp) && nic.GetIpType(c.srcIp) != "" {
		return errors.New("source ip and target ip must be the same type")
	}

	return cmdPing.CPingServer(c.srcIp, c.tarIp, c.count, c.size, c.interVal, c.timeOut, c.pattern)
}

func (c *cmdPingCommand) Commands() []simplecobra.Commander {
	return c.commands
}
