package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pgmtc/orchard-cli/internal/pkg/docker"
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
		"logs":    getComponentAction(runner.Logs),
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
