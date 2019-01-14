package local

import (
	"errors"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
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

func Test_status(t *testing.T) {
	common.SkipDockerTesting(t)
	COMPONENT_NAME := "test-status-component"
	IMAGE_NAME := "nginx:alpine"
	DOCKER_ID := "test-status-container"

	var config = common.MockConfig([]common.Component{
		common.Component{
			Name:          COMPONENT_NAME,
			Image:         IMAGE_NAME,
			DockerId:      DOCKER_ID,
			ContainerPort: 80,
			HostPort:      9998,
			TestUrl:       "http://localhost:9998",
		},
	})

	cmp := common.ComponentMap(config.CurrentProfile().Components)[COMPONENT_NAME]
	removeAction.Handler(common.ConsoleLogger{}, config, cmp) // Ignore error if it does not exist

	// Test - no image present
	if err := statusAction.Handler(common.ConsoleLogger{}, config); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Pull image and run status
	if err := pullAction.Handler(common.ConsoleLogger{}, config, cmp); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeAction.Handler(common.ConsoleLogger{}, config, cmp)
	if err := statusAction.Handler(common.ConsoleLogger{}, config); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Create container and run status
	if err := createAction.Handler(common.ConsoleLogger{}, config, cmp); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeAction.Handler(common.ConsoleLogger{}, config, cmp)
	if err := statusAction.Handler(common.ConsoleLogger{}, config); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Start container
	if err := startAction.Handler(common.ConsoleLogger{}, config, cmp); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := statusAction.Handler(common.ConsoleLogger{}, config); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Stop container
	if err := stopAction.Handler(common.ConsoleLogger{}, config, cmp); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := statusAction.Handler(common.ConsoleLogger{}, config); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Verbose and follow
	if err := statusAction.Handler(common.ConsoleLogger{}, config, "-v", "-f", "5"); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Follow only
	if err := statusAction.Handler(common.ConsoleLogger{}, config, "-f", "1"); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
}
