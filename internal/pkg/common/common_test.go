package common

import (
	"errors"
	"strings"
	"testing"
)

var (
	calledParam              string
	handlerMethodCalledStore = map[string]bool{}
)

// Method called by actionHandler test - success
func handlerMethod_success(component Component) error {
	handlerMethodCalledStore[component.Name] = true
	return nil
}

// Method called by actionHandler test - failure
func handlerMethod_fail(component Component) error {
	return errors.New("Method deliberately returned error")
}

func TestMakeActions(t *testing.T) {
	actions := MakeActions()
	if actions == nil {
		t.Fail()
	}
}

func TestMakeParams_helpAdded(t *testing.T) {
	actions := mockActions()
	args := []string{"help"}
	result := ParseParams(actions, args)
	if result != nil {
		t.Errorf("Expected help method to be added and call successful")
	}
	var providedFailingHelpHandler = func(args []string) error {
		return errors.New("deliberate error from provided help function")
	}
	actions["help"] = providedFailingHelpHandler
	result = ParseParams(actions, args)
	if result == nil {
		t.Errorf("Expected provided help method to be run and result in fail")
	}

}

func TestMakeParams_noArgs(t *testing.T) {
	actions := mockActions()
	var args []string
	result := ParseParams(actions, args)
	expectedMessagePrefix := "Action actionNonExisting' does not exist. Available actions ="

	if result == nil {
		t.Errorf("Expected method to fail (no parameters)")
	}
	if strings.HasPrefix(result.Error(), expectedMessagePrefix) {
		t.Errorf("Expected error message to start with '%s' but got '%s'", expectedMessagePrefix, result.Error())
	}
}

func TestMakeParams_success(t *testing.T) {
	actions := mockActions()
	args := []string{"actionSuccess", "param1"}
	result := ParseParams(actions, args)
	if result != nil {
		t.Errorf("Expected method to pass (return nil)")
	}
	if calledParam != "param1" {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestMakeParams_fail(t *testing.T) {
	actions := mockActions()
	args := []string{"actionFail", "param2"}
	result := ParseParams(actions, args)
	expectedMessage := "Action had failed"

	if result == nil {
		t.Errorf("Expected method to fail (error from action)")
	}
	if result.Error() != expectedMessage {
		t.Errorf("Expected error to be '%s' but got '%s'", expectedMessage, result.Error())
	}
	if calledParam != "param2" {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestMakeParams_nonExisting(t *testing.T) {
	actions := mockActions()
	args := []string{"actionNonExisting", "param3"}
	result := ParseParams(actions, args)
	expectedMessagePrefix := "Action actionNonExisting' does not exist. Available actions ="

	if result == nil {
		t.Errorf("Expected method to fail (non existing action)")
	}
	if strings.HasPrefix(result.Error(), expectedMessagePrefix) {
		t.Errorf("Expected error message to start with '%s' but got '%s'", expectedMessagePrefix, result.Error())
	}
	if calledParam != "param2" {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestGetActionNames(t *testing.T) {
	actionNames := getActionNames(mockActions())
	if len(actionNames) != 2 {
		t.Errorf("Unexpected number of action names: %d (expected 2)", len(actionNames))
	}
	if !ArrContains(actionNames, "actionSuccess") || !ArrContains(actionNames, "actionFail") {
		t.Errorf("Unexpected action names. Expected both actionSuccess and actionFail to be present, got %s", actionNames)
	}
}

func mockActions() map[string]func(args []string) error {
	actions := MakeActions()
	actions["actionSuccess"] = func(args []string) error {
		calledParam = args[0]
		return nil
	}
	actions["actionFail"] = func(args []string) error {
		calledParam = args[0]
		return errors.New("Action had failed")
	}
	return actions
}

func Test_componentActionHandler_missingArguments(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	actionHandler := ComponentActionHandler(handlerMethod_success)
	// Run without the arguments
	err := actionHandler([]string{}) // Pick one of the existing components
	if err == nil {
		t.Errorf("Expected error to be returned when no component provided")
	}
}

func Test_componentActionHandler_nonExistingComponent(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	actionHandler := ComponentActionHandler(handlerMethod_success)
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
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	actionHandler := ComponentActionHandler(handlerMethod_success)
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
	actionFailureHandler := ComponentActionHandler(handlerMethod_fail)
	err = actionFailureHandler([]string{"db"})
	if err == nil {
		t.Errorf("Expected error to be returned by failure handler method, but got no error")
	}
}

func Test_componentActionHandler_multiple(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	actionHandler := ComponentActionHandler(handlerMethod_success)
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
	actionFailureHandler := ComponentActionHandler(handlerMethod_fail)
	err = actionFailureHandler(validComponents)
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}

func Test_componentActionHandler_all(t *testing.T) {
	// Reset
	handlerMethodCalledStore = map[string]bool{}
	// Use ComponentActionHandler wrapper to convert it to actionHandler
	actionHandler := ComponentActionHandler(handlerMethod_success)
	// Run for all components
	err := actionHandler([]string{"all"}) // Pick one of the existing components
	if err != nil {
		t.Errorf("Expected no error to be returned, but got %s", err.Error())
	}
	// Check that handlerMethodCalled variable had been switched to true for all provided components
	componentsUnderTest := ComponentNames()
	for _, cmpUnderTest := range componentsUnderTest {
		if !handlerMethodCalledStore[cmpUnderTest] {
			t.Errorf("When using 'all' as a parameter, expected handlerMethod_success to be called for component %s, but it was not", cmpUnderTest)
		}
	}

	// Test failure scenario
	actionFailureHandler := ComponentActionHandler(handlerMethod_fail)
	// Run for all components
	err = actionFailureHandler([]string{"all"}) // Pick one of the existing components
	if err != nil {
		t.Errorf("Even though all runs for components have failed, expecting no error (those are reported as warning), but got %s", err.Error())
	}
}
