package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["build"] = common.ComponentActionHandler(build)
	return common.ParseParams(actions, args)
}
