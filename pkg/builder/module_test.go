package builder

import (
	"testing"
)

func TestModule_GetActions(t *testing.T) {
	mods := Module{}.GetActions()
	if len(mods) == 0 {
		t.Errorf("Unexpected number of modules: %d", len(mods))
	}
}
