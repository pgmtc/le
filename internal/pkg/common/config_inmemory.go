package common

type MemoryConfig struct {
	profile Profile
	config  Config
}

func (c MemoryConfig) CurrentProfileName() string {
	return "whatever"
}

func (MemoryConfig) SaveConfig() (fileName string, resultErr error) {
	return
}

func (MemoryConfig) LoadConfig() (resultErr error) {
	return
}

func (MemoryConfig) SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	return
}

func (MemoryConfig) LoadProfile(profileName string) (profile Profile, resultErr error) {
	return
}

func (c *MemoryConfig) CurrentProfile() Profile {
	return c.profile
}

func (MemoryConfig) GetAvailableProfiles() (profiles []string) {
	return []string{}
}

func (c *MemoryConfig) SetProfile(profileName string, profile Profile) {
	c.profile = profile
}

func (c *MemoryConfig) Config() Config {
	return c.config
}

func MockConfig(components []Component) Configuration {
	return &MemoryConfig{
		profile: Profile{
			Components: components,
		},
	}
}
