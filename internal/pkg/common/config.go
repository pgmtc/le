package common

var (
	CONFIG_LOCATION = "~/.orchard"
	CONFIG_FILE     = "config.yaml"

	CONFIG          Config
	CURRENT_PROFILE Profile
)

type Config struct {
	Profile string
}

func init() {
	CONFIG.Profile = "default"
	CURRENT_PROFILE = DefaultProfile()
}

func DefaultProfile() Profile {
	return Profile{
		SourceLocation: "",
		Components:     defaultComponents,
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
		Env:        []string{"SPRING_PROFILE=native,dcmp"},
		TestUrl:    "",
		DockerFile: "docker/Dockerfile-msvc",
		BuildRoot:  "modules/orchard-config-msvc",
	},
	Component{
		Name:     "auth",
		Image:    "orchard/orchard-auth-msvc:latest",
		DockerId: "dcmp_orchard-auth-msvc_1",
		Env:      []string{"SPRING_PROFILE=dcmp"},
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
		Env:      []string{"SPRING_PROFILE=dcmp"},
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
		Env:      []string{"SPRING_PROFILE=dcmp"},
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
		Env:           []string{"SPRING_PROFILE=dcmp"},
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
