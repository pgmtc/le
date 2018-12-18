package source

import (
	"github.com/pgmtc/orchard/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["pull"] = pull

	return common.ParseParams(actions, args)
}





