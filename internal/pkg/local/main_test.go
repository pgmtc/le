package local

import (
	"errors"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

var (
	handlerMethodCalledStore = map[string]bool{}
	handlerMethodLogsFollow  = false
)

// Method called by actionHandler test - success
func handlerMethod_success(component common.Component) error {
	handlerMethodCalledStore[component.Name] = true
	return nil
}

// Method called by actionHandler test - failure
func handlerMethod_fail(component common.Component) error {
	return errors.New("Method deliberately returned error")
}

// Method called by logsHandler test
func handlerMethod_logs(component common.Component, follow bool) error {
	handlerMethodLogsFollow = follow
	if follow { // Reuse follow flag for failure testing - lazy!
		return errors.New("deliberate error from logs handler")
	}
	return nil
}

func Test_componentActionHandler_missingArguments(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use componentActionHandler wrapper to convert it to actionHandler
	actionHandler := componentActionHandler(handlerMethod_success)
	// Run without the arguments
	err := actionHandler([]string{}) // Pick one of the existing components
	if err == nil {
		t.Errorf("Expected error to be returned when no component provided")
	}
}

func Test_componentActionHandler_nonExistingComponent(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use componentActionHandler wrapper to convert it to actionHandler
	actionHandler := componentActionHandler(handlerMethod_success)
	// Run the test for single component
	componentUnderTest := "nonexisting"                // Pick one of the existing components
	err := actionHandler([]string{componentUnderTest}) // Pick one of the existing components
	if err == nil {
		t.Errorf("Expected error to be returned for non-existing component")
	}
}

func Test_componentActionHandler_single(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use componentActionHandler wrapper to convert it to actionHandler
	actionHandler := componentActionHandler(handlerMethod_success)
	// Run the test for single component
	componentUnderTest := "db" // Pick one of the existing components
	err := actionHandler([]string{componentUnderTest})
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true
	if !handlerMethodCalledStore["db"] {
		t.Errorf("Expected handlerMethod_success to be called for component %s, but it was not", componentUnderTest)
	}
	// Run the test for non-existing component
	err = actionHandler([]string{"nonExisting"})
	if err == nil {
		t.Errorf("Expected error to be returned for non existing component")
	}
	// Test failure scenario
	actionFailureHandler := componentActionHandler(handlerMethod_fail)
	err = actionFailureHandler([]string{"db"})
	if err == nil {
		t.Errorf("Expected error to be returned by failure handler method, but got no error")
	}
}

func Test_componentActionHandler_multiple(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use componentActionHandler wrapper to convert it to actionHandler
	actionHandler := componentActionHandler(handlerMethod_success)
	// Run the test for single component
	validComponents := []string{"db", "case-flow", "auth"}
	allComponents := append(validComponents, "nonExisting")
	err := actionHandler(allComponents) // Pick one of the existing components
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true for all provided components
	for _, cmpUnderTest := range validComponents {
		if !handlerMethodCalledStore[cmpUnderTest] {
			t.Errorf("Expected handlerMethod_success to be called for component %s, but it was not", cmpUnderTest)
		}
	}
	// Check failure scenario
	actionFailureHandler := componentActionHandler(handlerMethod_fail)
	err = actionFailureHandler(validComponents)
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}

func Test_componentActionHandler_all(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use componentActionHandler wrapper to convert it to actionHandler
	actionHandler := componentActionHandler(handlerMethod_success)
	// Run for all components
	err := actionHandler([]string{"all"}) // Pick one of the existing components
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true for all provided components
	componentsUnderTest := common.ComponentNames()
	for _, cmpUnderTest := range componentsUnderTest {
		if !handlerMethodCalledStore[cmpUnderTest] {
			t.Errorf("When using 'all' as a parameter, expected handlerMethod_success to be called for component %s, but it was not", cmpUnderTest)
		}
	}

	// Test failure scenario
	actionFailureHandler := componentActionHandler(handlerMethod_fail)
	// Run for all components
	err = actionFailureHandler([]string{"all"}) // Pick one of the existing components
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
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
	// Not much testing here - just call it and hope for no error - also random depending on what is currently running or not - that's probably wrong
	err := status([]string{})
	if err != nil {
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
