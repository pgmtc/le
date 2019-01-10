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

func Test_logsHandler(t *testing.T) {
	// Reset for later usage
	handlerMethodLogsFollow = false

	// Test if handler is returned with follow = true
	logsActionHandlerFollowTrue := logsHandler(handlerMethod_logs, true)
	if logsActionHandlerFollowTrue == nil {
		t.Errorf("Expected logs action handler with follow = false to be created, but nil had been returned")
	}
	// Test if handler is returned with follow = false
	logsActionHandlerFollowFalse := logsHandler(handlerMethod_logs, false)
	if logsActionHandlerFollowFalse == nil {
		t.Errorf("Expected logs action handler with follow = true to be created, but nil had been returned")
	}

	var err error

	// Test logs coming back with error - (handlerMethod_logs reuses follow flag = true to fire an error)
	err = logsActionHandlerFollowTrue([]string{"db"})
	if err == nil {
		t.Errorf("Expected to fail, but no error has been retured")
	}
	// Test failure for non existing component
	err = logsActionHandlerFollowFalse([]string{"nonExisting"})
	if err == nil {
		t.Errorf("Expected to fail for nonExisting component, but no error has been retured")
	}
	// Test if follow had been correctly passed
	if handlerMethodLogsFollow == false {
		t.Errorf("Expected 'follow' flag = true to be passed to the logs handler, but it was not")
	}

	// Test logs coming back with success - (handlerMethod_logs reuses follow flag = false NOT to fire an error)
	err = logsActionHandlerFollowFalse([]string{"db"})
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Test if follow flag had been correctly passed
	if handlerMethodLogsFollow == true {
		t.Errorf("Expected 'follow' flag = false to be passed to the logs handler, but it was not")
	}

	// Test scenario when no component is passed
	err = logsActionHandlerFollowFalse([]string{})
	if err == nil {
		t.Errorf("Expected to fail when no component is passed")
	}

}

func Test_status(t *testing.T) {
	common.SkipDockerTesting(t)
	COMPONENT_NAME := "test-status-component"
	IMAGE_NAME := "nginx:alpine"
	DOCKER_ID := "test-status-container"

	common.CURRENT_PROFILE = common.Profile{
		Components: []common.Component{
			common.Component{
				Name:          COMPONENT_NAME,
				Image:         IMAGE_NAME,
				DockerId:      DOCKER_ID,
				ContainerPort: 80,
				HostPort:      9998,
				TestUrl:       "http://localhost:9998",
			},
		},
	}

	cmp := common.ComponentMap()[COMPONENT_NAME]
	removeImage(cmp, common.HandlerArguments{}) // Ignore error if it does not exist

	// Test - no image present
	if err := status([]string{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Pull image and run status
	if err := pullImage(cmp, common.HandlerArguments{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeImage(cmp, common.HandlerArguments{})
	if err := status([]string{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Create container and run status
	if err := createContainer(cmp, common.HandlerArguments{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	defer removeContainer(cmp, common.HandlerArguments{})
	if err := status([]string{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Start container
	if err := startContainer(cmp, common.HandlerArguments{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := status([]string{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

	// Stop container
	if err := stopContainer(cmp, common.HandlerArguments{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	if err := status([]string{}); err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}

}

func TestParse(t *testing.T) {
	// Prepare - stop redis container just in case it is running
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "failTest",
			args:    args{[]string{"nonExistingAction"}},
			wantErr: true,
		},
		{
			name:    "successTest",
			args:    args{[]string{"help"}},
			wantErr: false,
		}, // No need to test the rest - actions tested individually
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Parse(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
