package common

import (
	"reflect"
	"testing"
)

var mockConfig = MockConfig{}

func TestMockConfig_SaveConfig(t *testing.T) {
	_, err := mockConfig.SaveConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !mockConfig.saveConfigCalled {
		t.Error("saveConfigCalled expected to be true, is false")
	}

	mockConfig.reset().setSaveToFail()
	_, err = mockConfig.SaveConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if !mockConfig.saveConfigCalled {
		t.Error("saveConfigCalled expected to be true, is false")
	}
}

func TestMockConfig_LoadConfig(t *testing.T) {
	mockConfig.reset()
	err := mockConfig.LoadConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !mockConfig.loadConfigCalled {
		t.Error("loadConfigCalled expected to be true, is false")
	}

	mockConfig.reset().setLoadToFail()
	err = mockConfig.LoadConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if !mockConfig.loadConfigCalled {
		t.Error("loadConfigCalled expected to be true, is false")
	}
}

func TestMockConfig_SaveProfile(t *testing.T) {
	_, err := mockConfig.SaveProfile("test-profile", Profile{})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !mockConfig.saveProfileCalled {
		t.Error("saveProfileCalled expected to be true, is false")
	}

	mockConfig.reset().setSaveToFail()
	_, err = mockConfig.SaveProfile("test-profile", Profile{})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if !mockConfig.saveProfileCalled {
		t.Error("saveProfileCalled expected to be true, is false")
	}
}

func TestMockConfig_LoadProfile(t *testing.T) {
	mockConfig.reset()
	_, err := mockConfig.LoadProfile("test-profile")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if !mockConfig.loadProfileCalled {
		t.Error("loadProfileCalled expected to be true, is false")
	}

	mockConfig.reset().setLoadToFail()
	_, err = mockConfig.LoadProfile("test-profile")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if !mockConfig.loadProfileCalled {
		t.Error("loadProfileCalled expected to be true, is false")
	}
}

func TestMockConfig_Set_Current_Available_Profiles(t *testing.T) {
	profile := Profile{Components: []Component{{Name: "test-component"}}}
	mockConfig.SetProfile("test-profile", profile)
	if !mockConfig.setProfileCalled {
		t.Errorf("Expected setProfileCalled to be true")
	}

	// Current Profile Test
	if !reflect.DeepEqual(profile, mockConfig.CurrentProfile()) {
		t.Errorf("Expected profiles to match")
	}
	if !mockConfig.currentProfileCalled {
		t.Errorf("Expected currentProfileCalled to be true")
	}

	// GetAvailableProfiles test
	if len(mockConfig.GetAvailableProfiles()) != 1 {
		t.Errorf("Expected one profile to be available")
	}
	if mockConfig.GetAvailableProfiles()[0] != "test-profile" {
		t.Errorf("Expected available profile to equal test-profile, got %s", mockConfig.GetAvailableProfiles()[0])
	}
	if !mockConfig.getAvailableProfilesCalled {
		t.Errorf("Expected getAvailableProfielsCalled to be true")
	}

	// Test Reset
	mockConfig.reset()
	if mockConfig.failLoadRequired ||
		mockConfig.failSaveRequired ||
		mockConfig.configCalled ||
		mockConfig.setProfileCalled ||
		mockConfig.currentProfileCalled ||
		mockConfig.getAvailableProfilesCalled ||
		mockConfig.loadProfileCalled ||
		mockConfig.saveProfileCalled ||
		mockConfig.loadConfigCalled ||
		mockConfig.saveConfigCalled {
		t.Errorf("One of the flags is unexpectedly true")
	}
}

func TestSetBinLocation(t *testing.T) {
	mockConfig.setBinLocation("bin-location")
	if mockConfig.Config().BinLocation != "bin-location" {
		t.Errorf("Expected BinLocation to be %s, got %s", "bin-location", mockConfig.Config().BinLocation)
	}
}

func TestSetReleasesUrl(t *testing.T) {
	mockConfig.setReleasesUrl("releases-url")
	if mockConfig.Config().ReleasesURL != "releases-url" {
		t.Errorf("Expected ReleasesUrl to be %s, got %s", "releases-url", mockConfig.Config().ReleasesURL)
	}
}

func TestMockConfig_Config(t *testing.T) {
	mockConfig.Config()
	t.Skip()
}

func TestCreateMockConfig(t *testing.T) {
	components := []Component{
		{Name: "test-component"},
	}
	config := CreateMockConfig(components)
	if !reflect.DeepEqual(config.CurrentProfile().Components, components) {
		t.Errorf("Expected components to match provided")
	}

}
