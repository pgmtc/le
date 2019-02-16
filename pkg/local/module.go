package local

import (
	"github.com/pgmtc/le/pkg/common"
	"github.com/pgmtc/le/pkg/docker"
)

type Module struct{}

func (Module) GetActions() map[string]common.Action {
	runner := docker.DockerRunner{}
	return map[string]common.Action{
		"default": getRawAction(runner.Status),
		"status":  getRawAction(runner.Status),
		"create":  getComponentAction(runner.Create),
		"remove":  getComponentAction(runner.Remove),
		"start":   getComponentAction(runner.Start),
		"stop":    getComponentAction(runner.Stop),
		"pull":    getComponentAction(runner.Pull),
		"logs":    logsComponentAction(runner, false),
		"watch":   logsComponentAction(runner, true),
	}
}

func logsComponentAction(runner Runner, follow bool) common.Action {
	return &common.ComponentAction{
		Handler: func(ctx common.Context, cmp common.Component) error {
			return runner.Logs(ctx, cmp, follow)
		},
	}
}

func getComponentAction(handler func(ctx common.Context, cmp common.Component) error) common.Action {
	return &common.ComponentAction{
		Handler: handler,
	}
}

func getRawAction(handler func(ctx common.Context, args ...string) error) common.Action {
	return &common.RawAction{
		Handler: handler,
	}
}
