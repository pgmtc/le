package source

import (
	"testing"

	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var config = common.MockConfig([]common.Component{})

func Test_pullAction_Run(t *testing.T) {
	logger := common.ConsoleLogger{}
	err := pullAction{}.Run(logger, config)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
