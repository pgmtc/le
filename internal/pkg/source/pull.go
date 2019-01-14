package source

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var pullAction common.Action = &common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, cmp common.Component) error {
		log.Debugf("Pull source for %s\n", cmp.Name)
		return nil
	},
}
