package common

const CONFIG_FILE_NAME = "Config.yaml"

type Configuration interface {
	SaveConfig() (fileName string, resultErr error)
	LoadConfig() (resultErr error)
	SaveProfile(profileName string, profile Profile) (fileName string, resultErr error)
	LoadProfile(profileName string) (profile Profile, resultErr error)
	GetAvailableProfiles() (profiles []string)
	CurrentProfile() Profile
	SetProfile(profileName string, profile Profile)
	Config() Config
	SetRepositoryPrefix(url string)
}

type Config struct {
	Profile          string
	RepositoryPrefix string
}

var DefaultProfile = Profile{
	Components: defaultComponents,
}

type Profile struct {
	Components []Component
}

var defaultComponents = []Component{
	{
		Name:          "cmp1",
		Image:         "some-image-1:latest",
		ContainerPort: 8080,
		HostPort:      80,
		DockerId:      "container-1",
		TestUrl:       "http://localhost:80/",
		Env:           []string{"ENV_1=value", "ENV_2=value2"},
	},
	{
		Name:          "cmp2",
		Image:         "some-image-2:latest",
		DockerId:      "container-2",
		ContainerPort: 8443,
		HostPort:      443,
		TestUrl:       "http://localhost:443",
		Env:           []string{"ENV_1=value", "ENV_2=value2"},
		Links: []string{
			"container-1:cmp1",
		},
		Mounts: []string{"/host-dir:/container-dir"},
	},
}
