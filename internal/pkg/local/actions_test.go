package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func setUp() (config common.Configuration, log common.Logger, ctx common.Context) {
	config = common.CreateMockConfig([]common.Component{
		{
			Name:     "test-component",
			DockerId: "test-component",
			Image:    "bitnami/redis:latest",
		},
	})
	log = common.ConsoleLogger{}
	ctx = common.Context{
		Log:    log,
		Config: config,
	}
	return
}

func TestCreateAction(t *testing.T) {
	config, _, ctx := setUp()
	cmp := config.CurrentProfile().Components[0]
	err := createAction.Handler(ctx, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = startAction.Handler(ctx, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = stopAction.Handler(ctx, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = removeAction.Handler(ctx, cmp)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
