package source

import (
	"testing"

	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

var config = common.CreateMockConfig([]common.Component{
	{
		Name:  "testComponent",
		Image: "iron/go",
	},
})

func Test_pullAction_Run(t *testing.T) {
	logger := common.ConsoleLogger{}
	err := pullAction.Run(logger, config, "testComponent")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
