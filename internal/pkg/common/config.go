package common

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const CONFIG_FILE_NAME = "config.yaml"

type Configuration interface {
	SaveConfig() (fileName string, resultErr error)
	LoadConfig() (resultErr error)
	SaveProfile(profileName string, profile Profile) (fileName string, resultErr error)
	LoadProfile(profileName string) (profile Profile, resultErr error)
	CurrentProfile() Profile
	GetAvailableProfiles() (profiles []string)
}

type fileSystemConfig struct {
	configLocation string
	currentProfile Profile  // profile-xxx.yaml
	Config         struct { // config.yaml
		Profile     string
		ReleasesURL string
		BinLocation string
	}
}

func FileSystemConfig(configLocation string) (fcs Configuration, resultErr error) {
	fsc := &fileSystemConfig{
		configLocation: configLocation,
	}
	err := fsc.LoadConfig(CONFIG_FILE_NAME)
	if err != nil {
		resultErr = err
	}
	return
}

func (c *fileSystemConfig) LoadProfile(profileName string) (profile Profile, resultErr error) {
	configDir := c.initConfigDir(c.configLocation)
	out := Profile{}

	fileName := path.Join(configDir, "profile-"+profileName+".yaml")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		resultErr = errors.Errorf("profile does not exist, create it first")
		return
	}

	err := YamlUnmarshall(fileName, &out)
	if err != nil {
		resultErr = err
		return
	}

	profile = out
	return
}

func (c *fileSystemConfig) SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	configDir := c.initConfigDir(c.configLocation)
	fileName = path.Join(configDir, "profile-"+profileName+".yaml")
	if err := YamlMarshall(profile, fileName); err != nil {
		resultErr = errors.Errorf("error writing config file: %s", err.Error())
		return
	}
	return
}

func (c *fileSystemConfig) SaveConfig(configFileName string) (fileName string, resultErr error) {
	//_, err := c.SaveProfile(c.Config.Profile, c.CurrentProfile())
	//if err != nil {
	//	resultErr = errors.Errorf("error saving current profile: %s", err.Error())
	//	return
	//}
	fileName = path.Join(c.initConfigDir(c.configLocation), configFileName)
	if err := YamlMarshall(c.Config, fileName); err != nil {
		resultErr = errors.Errorf("error writing config file: %s", err.Error())
		return
	}
	return
}

func (c *fileSystemConfig) LoadConfig(configFileName string) (resultErr error) {
	//fileName := path.Join(c.initConfigDir(), configFileName)
	fileName := path.Join(c.configLocation, configFileName)
	if err := YamlUnmarshall(fileName, &c.Config); err != nil {
		resultErr = errors.Errorf("error reading config file %s: %s", fileName, err.Error())
		return
	}
	configProfile, err := c.LoadProfile(c.Config.Profile)
	if err != nil {
		resultErr = errors.Errorf("error loading config's profile: %s", err.Error())
		return
	}
	c.currentProfile = configProfile
	return
}

func (c *fileSystemConfig) GetAvailableProfiles() (profiles []string) {
	configDir := c.initConfigDir(c.configLocation)
	files, _ := filepath.Glob(configDir + "/profile-*.yaml")

	for _, file := range files {
		profileName := strings.TrimPrefix(file, configDir+"/profile-")
		profileName = strings.TrimSuffix(profileName, ".yaml")
		profiles = append(profiles, profileName)
	}
	return
}

func (c *fileSystemConfig) initConfigDir(configLocation string) (configDir string) {
	configDir = ParsePath(configLocation)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		color.Yellow("Config location '%s' does not exist, creating it", configDir)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return
}

func (c *fileSystemConfig) CurrentProfile() Profile {
	return c.currentProfile
}

var DefaultLocalProfile Profile = Profile{
	Components: defaultComponents,
}

var DefaultRemoteProfile Profile = Profile{
	Components: defaultRemoteComponents,
}

type Profile struct {
	Components []Component
}

var defaultComponents = []Component{
	{
		Name:          "db",
		Image:         "orchard/orchard-local-db:latest",
		ContainerPort: 3306,
		HostPort:      3306,
		DockerId:      "orchard-local-db",
		TestUrl:       "",
		DockerFile:    "modules/orchard-docker-local-db/Dockerfile",
		//DockerFile: "/Users/mfa/orchard/orchard-poc-umbrella/modules/orchard-docker-local-db/Dockerfile",
		BuildRoot: "modules/orchard-docker-local-db",
		//BuildRoot: "/Users/mfa/orchard/orchard-poc-umbrella/modules/orchard-docker-local-db",
	},
	{
		Name:     "redis",
		Image:    "bitnami/redis:latest",
		DockerId: "dcmp_orchard-redis_1",
		Env:      []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:  ""},
	{
		Name:       "config",
		Image:      "orchard/orchard-config-msvc:latest",
		DockerId:   "dcmp_orchard-config-msvc_1",
		Env:        []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:    "",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-config-msvc",
	},
	{
		Name:     "auth",
		Image:    "orchard/orchard-auth-msvc:latest",
		DockerId: "dcmp_orchard-auth-msvc_1",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-auth-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-auth-msvc",
	},
	{
		Name:     "doc-analysis",
		Image:    "orchard/orchard-doc-analysis-msvc:latest",
		DockerId: "dcmp_orchard-doc-analysis-msvc_1",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-doc-analysis-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-doc-analysis-msvc",
	},
	{
		Name:     "case-flow",
		Image:    "orchard/orchard-case-flow-msvc:latest",
		DockerId: "dcmp_orchard-case-flow-msvc_1",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
			"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-case-flow-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-case-flow-msvc",
	},
	{
		Name:          "gateway",
		Image:         "orchard/orchard-gateway-msvc:latest",
		DockerId:      "dcmp_orchard-gateway-msvc_1",
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		ContainerPort: 8080,
		HostPort:      8765,
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
			"dcmp_orchard-auth-msvc_1:auth",
			"dcmp_orchard-case-flow-msvc_1:case-flow",
			"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-gateway-msvc",
	},
	{
		Name:          "ui",
		Image:         "orchard/orchard-doc-analysis-ui:latest",
		DockerId:      "dcmp_orchard-doc-analysis-ui_1",
		ContainerPort: 80,
		HostPort:      3000,
		TestUrl:       "http://localhost:3000/",
		DockerFile:    "docker/Dockerfile-orchard-doc-analysis-ui",
		BuildRoot:     "modules/orchard-doc-analysis-ui/",
	},
}

var defaultRemoteComponents = []Component{
	{
		Name:          "db",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-local-db:latest",
		ContainerPort: 3306,
		HostPort:      3306,
		DockerId:      "orchard-local-db",
		TestUrl:       "",
	},
	{
		Name:          "redis",
		Image:         "bitnami/redis:latest",
		DockerId:      "dcmp_orchard-redis_1",
		ContainerPort: 6379,
		HostPort:      6379,
		Env:           []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:       "",
	},
	{
		Name:          "config",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-config-msvc:0.0.198",
		DockerId:      "dcmp_orchard-config-msvc_1",
		ContainerPort: 8080,
		HostPort:      8080,
		Env:           []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:       "http://localhost:8080/orchard-config-msvc/health",
	},
	{
		Name:          "auth",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-auth-msvc:0.0.164",
		DockerId:      "dcmp_orchard-auth-msvc_1",
		ContainerPort: 8080,
		HostPort:      50170,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
		},
		TestUrl: "http://localhost:50170/orchard-auth-msvc/health",
	},
	{
		Name:          "doc-analysis",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-doc-analysis-msvc:0.0.263",
		DockerId:      "dcmp_orchard-doc-analysis-msvc_1",
		ContainerPort: 8080,
		HostPort:      50130,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
		},
		TestUrl: "http://localhost:50130/orchard-doc-analysis-msvc/health",
	},
	{
		Name:          "case-flow",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-case-flow-msvc:0.0.323",
		DockerId:      "dcmp_orchard-case-flow-msvc_1",
		ContainerPort: 8080,
		HostPort:      50160,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
			"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
		},
		TestUrl: "http://localhost:50160/orchard-case-flow-msvc/health",
	},
	{
		Name:          "gateway",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-gateway-msvc:0.0.131",
		DockerId:      "dcmp_orchard-gateway-msvc_1",
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		ContainerPort: 8080,
		HostPort:      8765,
		Links: []string{
			"dcmp_orchard-redis_1:redis",
			"orchard-local-db:db",
			"dcmp_orchard-config-msvc_1:config",
			"dcmp_orchard-auth-msvc_1:auth",
			"dcmp_orchard-case-flow-msvc_1:case-flow",
			"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
		},
		TestUrl: "http://localhost:8765/orchard-gateway-msvc/health",
	},
	{
		Name:          "ui",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/temp-orchard-doc-analysis-ui:latest",
		DockerId:      "dcmp_orchard-doc-analysis-ui_1",
		ContainerPort: 80,
		HostPort:      3000,
		TestUrl:       "http://localhost:3000/",
	},
}
