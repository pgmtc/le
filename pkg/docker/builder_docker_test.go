package docker

import (
	"github.com/pgmtc/le/pkg/common"
	"testing"
)

func TestBuilder_BuildImage(t *testing.T) {
	// WIP:
	builder := Builder{}
	err := builder.BuildImage(common.Context{}, "", "", "", []string{}, false)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}
