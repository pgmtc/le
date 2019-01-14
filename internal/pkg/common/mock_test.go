package common

type mockConfig struct {
	profile Profile
}

func (mockConfig) SaveConfig() (fileName string, resultErr error) {
	return
}

func (mockConfig) LoadConfig() (resultErr error) {
	return
}

func (mockConfig) SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	return
}

func (mockConfig) LoadProfile(profileName string) (profile Profile, resultErr error) {
	return
}

func (c *mockConfig) CurrentProfile() Profile {
	return c.profile
}

func (mockConfig) GetAvailableProfiles() (profiles []string) {
	return []string{}
}

func MockConfig(components []Component) Configuration {
	return &mockConfig{
		profile: Profile{
			Components: components,
		},
	}

}
