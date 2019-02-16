package source

import (
	"github.com/pgmtc/le/pkg/common"
)

type Module struct{}

func (Module) GetActions() map[string]common.Action {
	return map[string]common.Action{
		"default": pullAction,
		"pull":    pullAction,
	}
}
