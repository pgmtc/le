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
}

type Config struct {
	Profile string
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
		DockerId:      "local_orchard-db",
		TestUrl:       "",
		DockerFile:    "modules/orchard-docker-local-db/Dockerfile",
		//DockerFile: "/Users/mfa/orchard/orchard-poc-umbrella/modules/orchard-docker-local-db/Dockerfile",
		BuildRoot: "modules/orchard-docker-local-db",
		//BuildRoot: "/Users/mfa/orchard/orchard-poc-umbrella/modules/orchard-docker-local-db",
	},
	{
		Name:     "redis",
		Image:    "bitnami/redis:latest",
		DockerId: "local_orchard-redis",
		Env:      []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:  ""},
	{
		Name:       "config",
		Image:      "orchard/orchard-config-msvc:latest",
		DockerId:   "local_orchard-config-msvc",
		Env:        []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:    "",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-config-msvc",
	},
	{
		Name:     "auth",
		Image:    "orchard/orchard-auth-msvc:latest",
		DockerId: "local_orchard-auth-msvc",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"local_orchard-redis:redis",
			"local_orchard-db:db",
			"local_orchard-config-msvc:config",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-auth-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-auth-msvc",
	},
	{
		Name:     "doc-analysis",
		Image:    "orchard/orchard-doc-analysis-msvc:latest",
		DockerId: "local_orchard-doc-analysis-msvc",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"local_orchard-redis:redis",
			"local_orchard-db:db",
			"local_orchard-config-msvc:config",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-doc-analysis-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-doc-analysis-msvc",
	},
	{
		Name:     "case-flow",
		Image:    "orchard/orchard-case-flow-msvc:latest",
		DockerId: "local_orchard-case-flow-msvc",
		Env:      []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"local_orchard-redis:redis",
			"local_orchard-db:db",
			"local_orchard-config-msvc:config",
			"local_orchard-doc-analysis-msvc:doc-analysis",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/orchard-case-flow-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-case-flow-msvc",
	},
	{
		Name:          "gateway",
		Image:         "orchard/orchard-gateway-msvc:latest",
		DockerId:      "local_orchard-gateway-msvc",
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		ContainerPort: 8080,
		HostPort:      8765,
		Links: []string{
			"local_orchard-redis:redis",
			"local_orchard-db:db",
			"local_orchard-config-msvc:config",
			"local_orchard-auth-msvc:auth",
			"local_orchard-case-flow-msvc:case-flow",
			"local_orchard-doc-analysis-msvc:doc-analysis",
		},
		TestUrl:    "http://localhost:8765/orchard-gateway-msvc/health",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-gateway-msvc",
	},
	{
		Name:          "ui",
		Image:         "orchard/orchard-doc-analysis-ui:latest",
		DockerId:      "local_orchard-doc-analysis-ui",
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
		DockerId:      "orchard-db",
		TestUrl:       "",
	},
	{
		Name:          "redis",
		Image:         "bitnami/redis:latest",
		DockerId:      "orchard-redis",
		ContainerPort: 6379,
		HostPort:      6379,
		Env:           []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:       "",
	},
	{
		Name:          "config",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-config-msvc:latest",
		DockerId:      "orchard-config-msvc",
		ContainerPort: 8080,
		HostPort:      8080,
		Env:           []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:       "http://localhost:8080/orchard-config-msvc/health",
	},
	{
		Name:          "auth",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-auth-msvc:latest",
		DockerId:      "orchard-auth-msvc",
		ContainerPort: 8080,
		HostPort:      50170,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"orchard-redis:redis",
			"orchard-db:db",
			"orchard-config-msvc:config",
		},
		TestUrl: "http://localhost:50170/orchard-auth-msvc/health",
	},
	{
		Name:          "doc-analysis",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-doc-analysis-msvc:latest",
		DockerId:      "orchard-doc-analysis-msvc",
		ContainerPort: 8080,
		HostPort:      50130,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"orchard-redis:redis",
			"orchard-db:db",
			"orchard-config-msvc:config",
		},
		TestUrl: "http://localhost:50130/orchard-doc-analysis-msvc/health",
	},
	{
		Name:          "case-flow",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-case-flow-msvc:latest",
		DockerId:      "orchard-case-flow-msvc",
		ContainerPort: 8080,
		HostPort:      50160,
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		Links: []string{
			"orchard-redis:redis",
			"orchard-db:db",
			"orchard-config-msvc:config",
			"orchard-doc-analysis-msvc:doc-analysis",
		},
		TestUrl: "http://localhost:50160/orchard-case-flow-msvc/health",
	},
	{
		Name:          "gateway",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-gateway-msvc:latest",
		DockerId:      "orchard-gateway-msvc",
		Env:           []string{"SPRING_PROFILES_ACTIVE=dcmp"},
		ContainerPort: 8080,
		HostPort:      8765,
		Links: []string{
			"orchard-redis:redis",
			"orchard-db:db",
			"orchard-config-msvc:config",
			"orchard-auth-msvc:auth",
			"orchard-case-flow-msvc:case-flow",
			"orchard-doc-analysis-msvc:doc-analysis",
		},
		TestUrl: "http://localhost:8765/orchard-gateway-msvc/health",
	},
	{
		Name:          "ui",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/temp-orchard-doc-analysis-ui:latest",
		DockerId:      "orchard-doc-analysis-ui",
		ContainerPort: 80,
		HostPort:      3000,
		TestUrl:       "http://localhost:3000/",
	},
}
