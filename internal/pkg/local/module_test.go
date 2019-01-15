package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func TestModule_GetActions(t *testing.T) {
	mods := Module{}.GetActions()
	if len(mods) == 0 {
		t.Errorf("Unexpected number of modules: %d", len(mods))
	}
}

func Test_getComponentAction(t *testing.T) {
	handlerRun := false
	handler := func(ctx common.Context, cmp common.Component) error {
		handlerRun = true
		return nil
	}

	action := getComponentAction(handler)
	action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{{Name: "test-component"}}),
	}, "test-component")
	if !handlerRun {
		t.Errorf("Handler has been expected to run, but it has not")
	}
}

func Test_getRawAction(t *testing.T) {
	handlerRun := false
	handler := func(ctx common.Context, args ...string) error {
		handlerRun = true
		return nil
	}

	action := getRawAction(handler)
	action.Run(common.Context{})
	if !handlerRun {
		t.Errorf("Handler has been expected to run, but it has not")
	}
}
