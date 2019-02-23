package common

import (
	"github.com/pkg/errors"
)

type MockConfig struct {
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
	currentProfile     Profile
	config             Config
}

func (c *MockConfig) SetRepositoryPrefix(url string) {
	c.config.RepositoryPrefix = url
}

func (c *MockConfig) SaveConfig(overwrite bool) (fileName string, resultErr error) {
	c.saveConfigCalled = true
	if c.failSaveRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *MockConfig) LoadConfig() (resultErr error) {
	c.loadConfigCalled = true
	if c.failLoadRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *MockConfig) SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	c.saveProfileCalled = true
	if c.failSaveRequired {
		resultErr = errors.New("Deliberate testing error")
	}
	return
}

func (c *MockConfig) LoadProfile(profileName string) (profile Profile, resultErr error) {
	c.loadProfileCalled = true
	if c.failLoadRequired {
		resultErr = errors.New("Deliberate testing error")
		return
	}
	profile = c.currentProfile
	return
}

func (c *MockConfig) GetAvailableProfiles() (profiles []string) {
	c.getAvailableProfilesCalled = true
	return []string{c.currentProfileName}
}

func (c *MockConfig) CurrentProfile() Profile {
	c.currentProfileCalled = true
	return c.currentProfile
}

func (c *MockConfig) SetProfile(profileName string, profile Profile) {
	c.setProfileCalled = true
	c.currentProfileName = profileName
	c.currentProfile = profile
}

func (c *MockConfig) Config() Config {
	c.configCalled = true
	return c.config
}

func (c *MockConfig) reset() *MockConfig {
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

func (c *MockConfig) setLoadToFail() *MockConfig {
	c.failLoadRequired = true
	return c
}

func (c *MockConfig) setSaveToFail() *MockConfig {
	c.failSaveRequired = true
	return c
}

func CreateMockConfig(components []Component) Configuration {
	config := MockConfig{
		currentProfile: Profile{
			Components: components,
		},
	}
	return &config
}
