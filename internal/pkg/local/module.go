package local

import "github.com/pgmtc/orchard-cli/internal/pkg/common"

type Module struct{}

func (Module) GetActions() map[string]common.Action {
	return map[string]common.Action{
		"default": &statusAction,
		"status":  &statusAction,
		"create":  &createAction,
		"remove":  &removeAction,
		"start":   &startAction,
		"stop":    &stopAction,
		"pull":    &pullAction,
		"logs":    &logsAction,
	}
}
