package common

import (
	"errors"
	"strings"
	"testing"
)

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
