package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["build"] = common.ComponentActionHandler(buildActionHandler, common.HandlerArguments{})
	return common.ParseParams(actions, args)
}

func buildActionHandler(cmp common.Component, arguments common.HandlerArguments) error {
	_, err := build(cmp, arguments)
	return err
}
