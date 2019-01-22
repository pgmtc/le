package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
	"path"
	"testing"
)

func cleanup() {
	if _, err := os.Stat(BUILDER_DIR); !os.IsNotExist(err) {
		os.RemoveAll(BUILDER_DIR)
	}
}

func Test(t *testing.T) {
	cleanup()
	defer cleanup()

	action := &initAction

	if err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	}); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	_, err := os.Stat(BUILDER_DIR)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if os.IsNotExist(err) {
		t.Errorf("%s directory had not been created", BUILDER_DIR)
	}

	// Check the contents
	expectedFiles := []string{"Dockerfile", "config.yaml"}
	for _, file := range expectedFiles {
		if _, err := os.Stat(path.Join(BUILDER_DIR, file)); os.IsNotExist(err) {
			t.Errorf("Expected file %s had not been found", file)
		}
	}

	// Test failure
	if err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	}); err == nil {
		t.Errorf("Expected error, got nothing")
	}
}
