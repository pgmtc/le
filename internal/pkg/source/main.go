package source

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["pull"] = pull

	return common.ParseParams(actions, args)

	// Get latest sources
	// Update
}





