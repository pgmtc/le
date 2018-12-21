package local

import (
	"github.com/pkg/errors"
	"testing"
)

var (
	handlerMethodCalledStore = map[string]bool{}
)

// Method called by actionHandler test - success
func handlerMethod_success(component Component) error {
	handlerMethodCalledStore[component.name] = true
	return nil
}

// Method called by actionHandler test - failure
func handlerMethod_fail(component Component) error {
	return errors.New("Method deliberately returned error")
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
	componentUnderTest := "nonexisting"                         // Pick one of the existing components
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
	err = actionFailureHandler([]string {"db"})
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
	err := actionHandler([]string {"all"}) // Pick one of the existing components
	if (err !=nil) {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true for all provided components
	componentsUnderTest := componentNames()
	for _, cmpUnderTest := range componentsUnderTest {
		if !handlerMethodCalledStore[cmpUnderTest] {
			t.Errorf("When using 'all' as a parameter, expected handlerMethod_success to be called for component %s, but it was not", cmpUnderTest)
		}
	}

	// Test failure scenario
	actionFailureHandler := componentActionHandler(handlerMethod_fail)
	// Run for all components
	err = actionFailureHandler([]string {"all"}) // Pick one of the existing components
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}
