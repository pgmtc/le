package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

type MockRunner struct{}

func (MockRunner) Create(ctx common.Context, cmp common.Component) error { return nil }
func (MockRunner) Remove(ctx common.Context, cmp common.Component) error { return nil }
func (MockRunner) Start(ctx common.Context, cmp common.Component) error  { return nil }
func (MockRunner) Stop(ctx common.Context, cmp common.Component) error   { return nil }
func (MockRunner) Pull(ctx common.Context, cmp common.Component) error   { return nil }
func (MockRunner) Logs(ctx common.Context, cmp common.Component) error   { return nil }
func (MockRunner) Status(ctx common.Context, args ...string) error       { return nil }

func setUp() (ctx common.Context, runner Runner) {
	config := common.CreateMockConfig([]common.Component{
		{
			Name:     "test-component",
			DockerId: "test-component",
			Image:    "bitnami/redis:latest",
		},
	})
	log := common.ConsoleLogger{}
	ctx = common.Context{
		Log:    log,
		Config: config,
	}

	runner = MockRunner{}

	return
}

func TestCreateAction(t *testing.T) {
	ctx, runner := setUp()
	createAction := getComponentAction(runner.Create)
	startAction := getComponentAction(runner.Start)
	stopAction := getComponentAction(runner.Stop)
	removeAction := getComponentAction(runner.Remove)

	err := createAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = startAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = stopAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	err = removeAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_status(t *testing.T) {
	ctx, runner := setUp()
	statusAction := getRawAction(runner.Status)
	pullAction := getComponentAction(runner.Pull)
	createAction := getComponentAction(runner.Create)
	startAction := getComponentAction(runner.Start)
	stopAction := getComponentAction(runner.Stop)
	removeAction := getComponentAction(runner.Remove)

	COMPONENT_NAME := "test-status-component"
	IMAGE_NAME := "nginx:alpine"
	DOCKER_ID := "test-status-container"

	var config = common.CreateMockConfig([]common.Component{
		common.Component{
			Name:          COMPONENT_NAME,
			Image:         IMAGE_NAME,
			DockerId:      DOCKER_ID,
			ContainerPort: 80,
			HostPort:      9998,
			TestUrl:       "http://localhost:9998",
		},
	})

	ctx = common.Context{
		Log:    common.ConsoleLogger{},
		Config: config,
	}

	removeAction.Run(ctx, COMPONENT_NAME) // Ignore error if it does not exist

	// Test - no image present
	if err := statusAction.Run(ctx); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Pull image and run status
	if err := pullAction.Run(ctx, COMPONENT_NAME); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeAction.Run(ctx, COMPONENT_NAME)
	if err := statusAction.Run(ctx); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Create container and run status
	if err := createAction.Run(ctx, COMPONENT_NAME); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeAction.Run(ctx, COMPONENT_NAME)
	if err := statusAction.Run(ctx); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Start container
	if err := startAction.Run(ctx, COMPONENT_NAME); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := statusAction.Run(ctx); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Stop container
	if err := stopAction.Run(ctx, COMPONENT_NAME); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := statusAction.Run(ctx); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Verbose and follow
	if err := statusAction.Run(ctx, "-v", "-f", "5"); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Follow only
	if err := statusAction.Run(ctx, "-f", "1"); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
}
