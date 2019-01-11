package common

var (
	CONFIG_LOCATION = "~/.orchard"
	CONFIG_FILE     = "config.yaml"

	CONFIG          Config
	CURRENT_PROFILE Profile
)

type Config struct {
	Profile     string
	ReleasesURL string
	BinLocation string
}

func init() {
	CONFIG.Profile = "default"
	CURRENT_PROFILE = DefaultRemoteProfile()
	CONFIG.ReleasesURL = "https://github.com/pgmtc/orchard-cli/releases/latest"
	CONFIG.BinLocation = "/usr/local/bin/orchard"
}

func DefaultLocalProfile() Profile {
	return Profile{
		SourceLocation: "",
		Components:     defaultComponents,
	}
}

func DefaultRemoteProfile() Profile {
	return Profile{
		SourceLocation: "",
		Components:     defaultRemoteComponents,
	}
}

type Profile struct {
	SourceLocation string
	Components     []Component
}

var defaultComponents []Component = []Component{
	Component{
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
	Component{
		Name:     "redis",
		Image:    "bitnami/redis:latest",
		DockerId: "dcmp_orchard-redis_1",
		Env:      []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:  ""},
	Component{
		Name:       "config",
		Image:      "orchard/orchard-config-msvc:latest",
		DockerId:   "dcmp_orchard-config-msvc_1",
		Env:        []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:    "",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-config-msvc",
	},
	Component{
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
	Component{
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
	Component{
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
	Component{
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
	Component{
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

var defaultRemoteComponents []Component = []Component{
	Component{
		Name:          "db",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-local-db:latest",
		ContainerPort: 3306,
		HostPort:      3306,
		DockerId:      "orchard-local-db",
		TestUrl:       "",
	},
	Component{
		Name:          "redis",
		Image:         "bitnami/redis:latest",
		DockerId:      "dcmp_orchard-redis_1",
		ContainerPort: 6379,
		HostPort:      6379,
		Env:           []string{"ALLOW_EMPTY_PASSWORD=yes"},
		TestUrl:       "",
	},
	Component{
		Name:          "config",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-config-msvc:0.0.198",
		DockerId:      "dcmp_orchard-config-msvc_1",
		ContainerPort: 8080,
		HostPort:      8080,
		Env:           []string{"SPRING_PROFILES_ACTIVE=native,dcmp"},
		TestUrl:       "http://localhost:8080/orchard-config-msvc/health",
	},
	Component{
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
	Component{
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
	Component{
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
	Component{
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
	Component{
		Name:          "ui",
		Image:         "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/temp-orchard-doc-analysis-ui:latest",
		DockerId:      "dcmp_orchard-doc-analysis-ui_1",
		ContainerPort: 80,
		HostPort:      3000,
		TestUrl:       "http://localhost:3000/",
	},
}
