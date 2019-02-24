package docker

import (
	"fmt"
	"github.com/pgmtc/le/pkg/common"
	"strconv"
	"time"
)

type Runner struct {
}

func (Runner) Status(ctx common.Context, args ...string) error {
	config := ctx.Config
	var verbose bool
	var follow bool
	var followLength int

	if len(args) > 0 && args[0] == "-v" || len(args) > 1 && args[1] == "-v" {
		verbose = true
	}

	// This could be improved - generalized
	if len(args) > 0 && args[0] == "-f" || len(args) > 1 && args[1] == "-f" {
		follow = true
		switch true {
		case len(args) > 1 && args[0] == "-f":
			i, err := strconv.Atoi(args[1])
			if err == nil {
				followLength = i
			}
		case len(args) > 2 && args[1] == "-f":
			i, err := strconv.Atoi(args[2])
			if err == nil {
				followLength = i
			}
		}
		follow = true
	}

	if !follow {
		return printStatus(config.CurrentProfile().Components, verbose, follow, ctx.Log)
	}
	counter := 0
	for {
		printStatus(config.CurrentProfile().Components, verbose, follow, ctx.Log)
		fmt.Println("local status: ", time.Now().Format("2006-01-02 15:04:05"))
		counter++
		time.Sleep(1 * time.Second)
		if counter == followLength {
			break
		}
	}

	return nil
}

func (Runner) Create(ctx common.Context, cmp common.Component) error {
	return createContainer(cmp, ctx.Log.Infof)
}

func (Runner) Remove(ctx common.Context, cmp common.Component) error {
	return removeComponent(cmp, ctx.Log.Infof)
}

func (Runner) Start(ctx common.Context, cmp common.Component) error {
	return startComponent(cmp, ctx.Log.Infof)
}

func (Runner) Stop(ctx common.Context, cmp common.Component) error {
	return stopContainer(cmp, ctx.Log.Infof)
}

func (Runner) Pull(ctx common.Context, cmp common.Component) error {
	return pullImage(cmp, ctx.Log.Infof)
}

func (Runner) Logs(ctx common.Context, cmp common.Component, follow bool) error {
	return printLogs(cmp, follow)
}
