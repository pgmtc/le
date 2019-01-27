package repo

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var urlAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		var repoName string = ""
		if len(args) > 0 {
			repoName = args[0]
		}
		ctx.Log.Infof("git clone %s\n", ctx.Config.Config().RepositoryPrefix+repoName)
		return nil
	},
}
