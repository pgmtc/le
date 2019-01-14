package common

import (
	"testing"
)

var config = MemoryConfig{}

func TestMemoryConfig_CurrentProfileName(t *testing.T) {
	config.CurrentProfile()
	t.Skip()
}

func TestMemoryConfig_SaveConfig(t *testing.T) {
	config.SaveConfig()
	t.Skip()
}

func TestMemoryConfig_LoadConfig(t *testing.T) {
	config.LoadConfig()
	t.Skip()
}

func TestMemoryConfig_SaveProfile(t *testing.T) {
	config.SaveProfile("test", Profile{})
	t.Skip()
}

func TestMemoryConfig_LoadProfile(t *testing.T) {
	config.LoadProfile("test")
	t.Skip()
}

func TestMemoryConfig_CurrentProfile(t *testing.T) {
	config.CurrentProfile()
	t.Skip()
}

func TestMemoryConfig_GetAvailableProfiles(t *testing.T) {
	config.GetAvailableProfiles()
	t.Skip()
}

func TestMemoryConfig_SetProfile(t *testing.T) {
	config.SetProfile("test", Profile{})
	t.Skip()
}

func TestMemoryConfig_Config(t *testing.T) {
	config.Config()
	t.Skip()
}

func TestMockConfig(t *testing.T) {
	MockConfig([]Component{})
	t.Skip()
}
