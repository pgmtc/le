package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func setUp() (config common.Configuration, log common.Logger) {
	config = common.CreateMockConfig([]common.Component{
		{
			Name:     "test-component",
			DockerId: "test-component",
			Image:    "bitnami/redis:latest",
		},
	})
	log = common.ConsoleLogger{}
	return
}

func TestCreateAction(t *testing.T) {
	config, log := setUp()
	cmp := config.CurrentProfile().Components[0]
	err := createAction.Handler(log, config, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = startAction.Handler(log, config, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = stopAction.Handler(log, config, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = removeAction.Handler(log, config, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
