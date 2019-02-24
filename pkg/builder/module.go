package builder

import (
	"github.com/pgmtc/le/pkg/common"
	"github.com/pgmtc/le/pkg/docker"
)

type Module struct{}

func (Module) GetActions() map[string]common.Action {
	builder := docker.Builder{}
	return map[string]common.Action{
		"build": getBuildAction(builder),
		"init":  &initAction,
	}
}
