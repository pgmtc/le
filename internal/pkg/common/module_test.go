package common

import "testing"

func TestRawAction_Run(t *testing.T) {
	var handlerCalled = false
	var action = RawAction{
		Handler: func(log Logger, config Configuration, args ...string) error {
			handlerCalled = true
			return nil
		},
	}
	err := action.Run(ConsoleLogger{}, MockConfig([]Component{}))
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if !handlerCalled {
		t.Errorf("Expected handler to be called and swap false to true")
	}

}
