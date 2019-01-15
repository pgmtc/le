package local

import (
	"errors"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var (
	handlerMethodLogsFollow = false
)

// Method called by logsHandler test
func handlerMethod_logs(component common.Component, follow bool) error {
	handlerMethodLogsFollow = follow
	if follow { // Reuse follow flag for failure testing - lazy!
		return errors.New("deliberate error from logs handler")
	}
	return nil
}
