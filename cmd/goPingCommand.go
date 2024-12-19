/**
 * Package cmd
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/16 16:36
 */

package cmd

import (
	"context"
	"errors"
	"searchlight/internal/goPing"
	"searchlight/pkg/nic"
	"searchlight/pkg/simplecobra"
)

type (
	goPingCommand struct {
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
	}
)

func newGoPingCommand() *goPingCommand {
	c := &goPingCommand{}
	return c
}

func (g *goPingCommand) Name() string {
	return "goPing"
}

func (g *goPingCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "ping detection using the pro-bing package"
	cmd.Long = "Using p8s open source pro-ping bing scheme (https://github.com/prometheus-community/pro-bing) to realize network detection function. On this basis, we optimize the filtering of network packets and filter the required data at the kernel layer. Reduces cpu usage when processing packets"
	cmd.Aliases = []string{"gp", "gping", "gPing", "goping"}

	cmd.Flags().StringVarP(&g.srcIp, "srcIP", "I", "", "The detect of source ip or local ip")
	cmd.Flags().StringVarP(&g.tarIp, "targetIP", "T", "", "The detect of target ip or remote ip")
	cmd.Flags().IntVarP(&g.count, "count", "c", 8, "Set the number of times a request response is completed,default 8")
	cmd.Flags().IntVarP(&g.size, "size", "s", 24, "Set the packet size,default 24 bytes")
	cmd.Flags().IntVarP(&g.interVal, "interval", "i", 1, "Set the interval time,default 1 second")
	cmd.Flags().IntVarP(&g.timeOut, "timeout", "t", 30, "Set the timeout,default 30 second")
	return nil
}

func (g *goPingCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

// Run 具体的执行逻辑
func (g *goPingCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	// Check data
	if g.tarIp == "" {
		return errors.New("target ip is empty, please set it. e.g. -T 8.8.8.8")
	}

	if nic.GetIpType(g.srcIp) != nic.GetIpType(g.tarIp) && nic.GetIpType(g.srcIp) != "" {
		return errors.New("source ip and target ip must be the same type")
	}

	return goPing.GPService(g.srcIp, g.tarIp, g.count, g.size, g.interVal, g.timeOut)
}

func (g *goPingCommand) Commands() []simplecobra.Commander {
	return g.commands
}
