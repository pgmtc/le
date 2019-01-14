package config

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func setUp() (config *DummyConfig, log *common.StringLogger) {
	config = &DummyConfig{
		currentProfile: common.Profile{
			Components: []common.Component{
				{
					Name: "test-component",
				},
			},
		},
	}
	log = &common.StringLogger{}
	return
}

func TestCreateAction(t *testing.T) {
	config, logger := setUp()
	// Fail on missing parameter
	config.reset()
	if err := createAction.Handler(logger, config); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Fail in underlying config methods - one parameter
	config.reset().setSaveToFail()
	if err := createAction.Handler(logger, config, "sourceprofile"); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Fail in underlying config methods - two parameters
	config.reset().setLoadToFail()
	if err := createAction.Handler(logger, config, "sourceprofile", "destinationprofile"); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Success - one parameter
	config.reset()
	if err := createAction.Handler(logger, config, "sourceprofile"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !config.saveProfileCalled {
		t.Errorf("config.SaveProfile profile had not been called")
	}

	// Success - two parameters
	config.reset()
	if err := createAction.Handler(logger, config, "sourceprofile", "destinationprofile"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !config.saveProfileCalled {
		t.Errorf("config.SaveProfile profile had not been called")
	}

}

func TestInitAction(t *testing.T) {
	config, logger := setUp()
	if err := initAction.Handler(logger, config); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !config.saveConfigCalled {
		t.Errorf("Save config expected to have been called, it had not")
	}
	if !config.saveProfileCalled {
		t.Errorf("Save profile expected to have been called, it had not")
	}
}

func TestStatusAction(t *testing.T) {
	config, logger := setUp()
	logger = &common.StringLogger{}
	// Normal
	if err := statusAction.Handler(logger, config); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(logger.InfoMessages) != 4 {
		t.Errorf("Expected 4 info messages, got %d", len(logger.InfoMessages))
		t.Errorf(strings.Join(logger.InfoMessages, ""))
	}

	// Verbose
	if err := statusAction.Handler(logger, config, "-v"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(logger.InfoMessages) == 0 {
		t.Errorf("Nothing in logger info slice - that's unexpected for status")
	}
}

func TestSwitchAction(t *testing.T) {
	// Missing parameters
	config, logger := setUp()
	if err := switchAction.Handler(logger, config); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test of failures
	config.reset().setLoadToFail()
	if err := switchAction.Handler(logger, config, "new-profile"); err == nil {
		t.Errorf("Expected error, got nothing")
	}
	config.reset().setSaveToFail()
	if err := switchAction.Handler(logger, config, "new-profile"); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test of success
	config.reset()
	if err := switchAction.Handler(logger, config, "new-profile"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !config.loadProfileCalled {
		t.Errorf("Expected config's loadProfile to be called but it was not")
	}
	if !config.setProfileCalled {
		t.Errorf("Expected config's setProfile to be called but it was not")
	}
	if !config.saveConfigCalled {
		t.Errorf("Expected config's saveProfile to be called but it was not")
	}
}

func TestUpdateCli(t *testing.T) {
	config, logger := setUp()
	tmpDir, err := ioutil.TempDir("", "orchard-test")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir) // clean up

	config.setBinLocation(path.Join(tmpDir, "orchard-updated")).setReleasesUrl("https://github.com/pgmtc/orchard-cli/releases/latest")
	err = updateCliAction.Handler(logger, config)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
