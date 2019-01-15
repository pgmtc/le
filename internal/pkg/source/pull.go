package source

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var pullAcrtion common.Action = &common.ComponentAction{}
var pullAction common.Action = &common.ComponentAction{

	Handler: func(ctx common.Context, cmp common.Component) error {
		ctx.Log.Debugf("Pull source for %s\n", cmp.Name)
		return nil
	},
}
