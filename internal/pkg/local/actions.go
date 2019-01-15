package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var runner Runner = DockerRunner{}

var createAction = common.ComponentAction{
	Handler: runner.Create,
}

var removeAction = common.ComponentAction{
	Handler: runner.Remove,
}

var startAction = common.ComponentAction{
	Handler: runner.Start,
}

var stopAction = common.ComponentAction{
	Handler: runner.Stop,
}

var pullAction = common.ComponentAction{
	Handler: runner.Pull,
}

var logsAction = common.ComponentAction{
	Handler: runner.Logs,
}
