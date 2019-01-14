package common

import (
	"github.com/pkg/errors"
	"testing"
)

var (
	calledParam              string
	handlerMethodCalledStore = map[string]bool{}
	testConfig               = MockConfig([]Component{
		{Name: "test-component-1", DockerId: "test-component-1-docker-id", Image: "test-component-1-image"},
		{Name: "test-component-2", DockerId: "test-component-2-docker-id", Image: "test-component-2-image"},
	})
)

func actionHandlerMethod_success(log Logger, config Configuration, component Component) error {
	handlerMethodCalledStore[component.Name] = true
	return nil
}

// Method called by actionHandler test - failure
func actionHandlerMethod_fail(log Logger, config Configuration, component Component) error {
	return errors.New("Method deliberately returned error")
}

func Test_componentAction_missingArguments(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	action := ComponentAction{
		Handler: actionHandlerMethod_success,
	}

	err := action.Run(ConsoleLogger{}, testConfig)
	if err == nil {
		t.Errorf("Expected error to be returned when no component provided")
	}
}

func Test_componentAction_nonExistingComponent(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	action := ComponentAction{
		Handler: actionHandlerMethod_success,
	}
	// Run the test for single component
	componentUnderTest := "nonexisting" // Pick one of the existing components
	err := action.Run(ConsoleLogger{}, testConfig, componentUnderTest)
	if err == nil {
		t.Errorf("Expected error to be returned for non-existing component")
	}
}

func Test_componentAction_single(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	action := ComponentAction{
		Handler: actionHandlerMethod_success,
	}
	// Run the test for single component
	componentUnderTest := "test-component-1" // Pick one of the existing components
	err := action.Run(ConsoleLogger{}, testConfig, componentUnderTest)
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true
	if !handlerMethodCalledStore[componentUnderTest] {
		t.Errorf("Expected handlerMethod_success to be called for component %s, but it was not", componentUnderTest)
	}
	// Run the test for non-existing component
	err = action.Run(ConsoleLogger{}, testConfig, "nonExisting")
	if err == nil {
		t.Errorf("Expected error to be returned for non existing component")
	}
	// Test failure scenario
	actionFailure := ComponentAction{
		Handler: actionHandlerMethod_fail,
	}
	err = actionFailure.Run(ConsoleLogger{}, testConfig, "test-component-1")
	if err == nil {
		t.Errorf("Expected error to be returned by failure handler method, but got no error")
	}
}

func Test_componentActionHandler_multiple(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	action := ComponentAction{
		Handler: actionHandlerMethod_success,
	}

	// Run the test for single component
	validComponents := []string{"test-component-1", "test-component-2"}
	allComponents := append(validComponents, "nonExisting")
	err := action.Run(ConsoleLogger{}, testConfig, allComponents...) // Pick one of the existing components
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
	actionFailure := ComponentAction{
		Handler: actionHandlerMethod_fail,
	}
	err = actionFailure.Run(ConsoleLogger{}, testConfig, validComponents...)
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}

func Test_componentActionHandler_all(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	action := ComponentAction{
		Handler: actionHandlerMethod_success,
	}
	// Run for all components
	err := action.Run(ConsoleLogger{}, testConfig, "all") // Pick one of the existing components
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true for all provided components
	componentsUnderTest := ComponentNames(testConfig.CurrentProfile().Components)
	for _, cmpUnderTest := range componentsUnderTest {
		if !handlerMethodCalledStore[cmpUnderTest] {
			t.Errorf("When using 'all' as a parameter, expected handlerMethod_success to be called for component %s, but it was not", cmpUnderTest)
		}
	}

	// Test failure scenario
	actionFailure := ComponentAction{
		Handler: actionHandlerMethod_fail,
	}
	// Run for all components
	err = actionFailure.Run(ConsoleLogger{}, testConfig, "all")
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}
