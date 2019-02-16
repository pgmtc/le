package config

import (
	"github.com/pgmtc/le/pkg/common"
	"github.com/pkg/errors"
)

type DummyConfig struct {
	failSaveRequired           bool
	failLoadRequired           bool
	saveConfigCalled           bool
	loadConfigCalled           bool
	saveProfileCalled          bool
	loadProfileCalled          bool
	getAvailableProfilesCalled bool
	currentProfileCalled       bool
	setProfileCalled           bool
	configCalled               bool

	currentProfileName string
	currentProfile     common.Profile
	config             common.Config
}

func (c *DummyConfig) SaveConfig() (fileName string, resultErr error) {
	c.saveConfigCalled = true
	if c.failSaveRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *DummyConfig) LoadConfig() (resultErr error) {
	c.loadConfigCalled = true
	if c.failLoadRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *DummyConfig) SaveProfile(profileName string, profile common.Profile) (fileName string, resultErr error) {
	c.saveProfileCalled = true
	if c.failSaveRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *DummyConfig) LoadProfile(profileName string) (profile common.Profile, resultErr error) {
	c.loadProfileCalled = true
	if c.failLoadRequired {
		resultErr = errors.New("Deliberate testing error")
		return
	}
	profile = c.currentProfile
	return
}

func (c *DummyConfig) GetAvailableProfiles() (profiles []string) {
	c.getAvailableProfilesCalled = true
	return []string{c.currentProfileName}
}

func (c *DummyConfig) CurrentProfile() common.Profile {
	c.currentProfileCalled = true
	return c.currentProfile
}

func (c *DummyConfig) SetProfile(profileName string, profile common.Profile) {
	c.setProfileCalled = true
	c.currentProfileName = profileName
	c.currentProfile = profile
}

func (c *DummyConfig) Config() common.Config {
	c.configCalled = true
	return c.config
}

func (c *DummyConfig) reset() *DummyConfig {
	c.failLoadRequired = false
	c.failSaveRequired = false
	c.configCalled = false
	c.setProfileCalled = false
	c.currentProfileCalled = false
	c.getAvailableProfilesCalled = false
	c.loadProfileCalled = false
	c.saveProfileCalled = false
	c.loadConfigCalled = false
	c.saveConfigCalled = false
	return c
}

func (c *DummyConfig) setLoadToFail() *DummyConfig {
	c.failLoadRequired = true
	return c
}

func (c *DummyConfig) setSaveToFail() *DummyConfig {
	c.failSaveRequired = true
	return c
}
