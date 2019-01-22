package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

type Module struct{}

func (Module) GetActions() map[string]common.Action {
	return map[string]common.Action{
		"build": &buildAction,
		"init":  &initAction,
	}
}
