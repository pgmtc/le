package common

import (
	"github.com/pkg/errors"
	"testing"
)

var (
	calledParam string
)

func TestMakeActions(t *testing.T) {
	actions := MakeActions()
	if (actions == nil) {
		t.Fail()
	}
}

func TestMakeParams_noArgs(t *testing.T) {
	actions := mockActions()
	var args []string
	result := ParseParams(actions, args)
	expectedMessage1 := "Missing action, available actions = [actionSuccess actionFail]"
	expectedMessage2 := "Missing action, available actions = [actionFail actionSuccess]"

	if (result == nil) {
		t.Errorf("Expected method to fail (no parameters)")
	}
	if (result.Error() != expectedMessage1 && result.Error() != expectedMessage2) {
		t.Errorf("Expected error to be '%s' but got '%s'", expectedMessage1, result.Error())
	}
}

func TestMakeParams_success(t *testing.T) {
	actions := mockActions()
	args := []string{"actionSuccess", "param1"}
	result := ParseParams(actions, args)
	if (result != nil) {
		t.Errorf("Expected method to pass (return nil)")
	}
	if (calledParam != "param1") {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestMakeParams_fail(t *testing.T) {
	actions := mockActions()
	args := []string{"actionFail", "param2"}
	result := ParseParams(actions, args)
	expectedMessage := "Action had failed"

	if (result == nil) {
		t.Errorf("Expected method to fail (error from action)")
	}
	if (result.Error() != expectedMessage) {
		t.Errorf("Expected error to be '%s' but got '%s'", expectedMessage, result.Error())
	}
	if (calledParam != "param2") {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestMakeParams_nonExisting(t *testing.T) {
	actions := mockActions()
	args := []string{"actionNonExisting", "param3"}
	result := ParseParams(actions, args)
	expectedMessage1 := "Action 'actionNonExisting' does not exist. Available actions = [actionSuccess actionFail]"
	expectedMessage2 := "Action 'actionNonExisting' does not exist. Available actions = [actionSuccess actionFail]"

	if (result == nil) {
		t.Errorf("Expected method to fail (non existing action)")
	}
	if (result.Error() != expectedMessage1 && result.Error() != expectedMessage2) {
		t.Errorf("Expected error to be '%s' but got '%s'", expectedMessage1, result.Error())
	}
	if (calledParam != "param2") {
		t.Errorf("Expected that method receives right parameter value")
	}
}

func TestGetActionNames(t *testing.T) {
	actionNames := getActionNames(mockActions())
	if (len(actionNames) != 2) {
		t.Errorf("Unexpected number of action names: %d (expected 2)", len(actionNames))
	}
	if (actionNames[0] != "actionSuccess" || actionNames[1] != "actionFail") {
		t.Errorf("Unexpected action names. Expected [actionSuccess, actionFail], got %s", actionNames)
	}
}


func mockActions() map[string]func(args []string) error {
	actions := MakeActions()
	actions["actionSuccess"] = func (args[] string) error {
		calledParam = args[0]
		return nil
	}
	actions["actionFail"] = func (args[] string) error {
		calledParam = args[0]
		return errors.New("Action had failed")
	}
	return actions
}

