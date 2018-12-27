package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func Test_build(t *testing.T) {
	var err error
	cmp := common.ComponentMap()["db"]
	err = build(cmp)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
}
