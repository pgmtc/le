package common

type Component struct {
	Name          string
	DockerId      string
	TestUrl       string
	Image         string
	ContainerPort int
	HostPort      int
	Env           []string
	Links         []string
	Volumes       []string
}

func ComponentNames() []string {
	components := GetComponents()
	componentNames := []string{}
	for _, component := range components {
		componentNames = append(componentNames, component.Name)
	}
	return componentNames
}

func ComponentMap() map[string]Component {
	components := GetComponents()
	elementMap := make(map[string]Component)
	for i := 0; i < len(components); i += 1 {
		elementMap[components[i].Name] = components[i]
	}
	return elementMap
}

func GetComponents() []Component {
	return []Component{
		Component{
			Name:          "db",
			Image:         "orchard/orchard-local-db:latest",
			ContainerPort: 3306,
			HostPort:      3306,
			DockerId:      "orchard-local-db",
			TestUrl:       ""},
		Component{
			Name:     "redis",
			Image:    "bitnami/redis:latest",
			DockerId: "dcmp_orchard-redis_1",
			Env:      []string{"ALLOW_EMPTY_PASSWORD=yes"},
			TestUrl:  ""},
		Component{
			Name:     "config",
			Image:    "orchard/orchard-config-msvc:latest",
			DockerId: "dcmp_orchard-config-msvc_1",
			TestUrl:  ""},
		Component{
			Name:     "auth",
			Image:    "orchard/orchard-auth-msvc:latest",
			DockerId: "dcmp_orchard-auth-msvc_1",
			Links: []string{
				"dcmp_orchard-redis_1:redis",
				"orchard-local-db:db",
				"dcmp_orchard-config-msvc_1:config",
			},
			TestUrl: "http://localhost:8765/orchard-gateway-msvc/orchard-auth-msvc/health"},
		Component{
			Name:     "doc-analysis",
			Image:    "orchard/orchard-doc-analysis-msvc:latest",
			DockerId: "dcmp_orchard-doc-analysis-msvc_1",
			Links: []string{
				"dcmp_orchard-redis_1:redis",
				"orchard-local-db:db",
				"dcmp_orchard-config-msvc_1:config",
			},
			TestUrl: "http://localhost:8765/orchard-gateway-msvc/orchard-doc-analysis-msvc/health"},
		Component{
			Name:     "case-flow",
			Image:    "orchard/orchard-case-flow-msvc:latest",
			DockerId: "dcmp_orchard-case-flow-msvc_1",
			Links: []string{
				"dcmp_orchard-redis_1:redis",
				"orchard-local-db:db",
				"dcmp_orchard-config-msvc_1:config",
				"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
			},
			TestUrl: "http://localhost:8765/orchard-gateway-msvc/orchard-case-flow-msvc/health"},
		Component{
			Name:          "gateway",
			Image:         "orchard/orchard-gateway-msvc:latest",
			DockerId:      "dcmp_orchard-gateway-msvc_1",
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
			TestUrl: "http://localhost:8765/orchard-gateway-msvc/health"},
		Component{
			Name:          "ui",
			Image:         "orchard/orchard-doc-analysis-ui:latest",
			DockerId:      "dcmp_orchard-doc-analysis-ui_1",
			ContainerPort: 80,
			HostPort:      3000,
			TestUrl:       "http://localhost:3000/"},
	}
}
