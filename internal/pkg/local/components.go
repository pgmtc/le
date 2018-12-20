package local

type Component struct {
	name     string
	dockerId string
	testUrl  string
	image string
	containerPort int
	hostPort int
	env []string
	links []string
	volumes []string
}

func componentNames() []string {
	components := getComponents()
	componentNames := []string{}
	for _, component := range components {
		componentNames = append(componentNames, component.name)
	}
	return componentNames
}

func componentMap() map[string]Component {
	components := getComponents()
	elementMap := make(map[string]Component)
	for i := 0; i < len(components); i += 1 {
		elementMap[components[i].name] = components[i]
	}
	return elementMap
}

func getComponents() []Component {
	return []Component{
		Component{
			name: "db",
			image: "orchard/orchard-local-db:latest",
			containerPort: 3306,
			hostPort: 3306,
			dockerId: "orchard-local-db",
			testUrl:  ""},
		Component{
			name: "redis",
			image: "bitnami/redis:latest",
			dockerId: "dcmp_orchard-redis_1",
			env: []string {"ALLOW_EMPTY_PASSWORD=yes"},
			testUrl:  ""},
		Component{
			name: "config",
			image: "orchard/orchard-config-msvc:latest",
			dockerId: "dcmp_orchard-config-msvc_1",
			testUrl:  ""},
		Component{
			name: "gateway",
			image: "orchard/orchard-gateway-msvc:latest",
			dockerId: "dcmp_orchard-gateway-msvc_1",
			containerPort: 8080,
			hostPort: 8765,
			links: []string {
				"dcmp_orchard-redis_1:redis",
				"dcmp_orchard-config-msvc_1:config",
				"dcmp_orchard-auth-msvc_1:auth",
				"dcmp_orchard-case-flow_1:case-flow",
				"dcmp_orchard-doc-analysis-msvc_1:doc-analysis",
			},
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/health"},

		Component{
			name: "auth",
			image: "orchard/orchard-auth-msvc:latest",
			dockerId: "dcmp_orchard-auth-msvc_1",
			links: []string {
				"dcmp_orchard-redis_1:redis",
				"dcmp_orchard-config-msvc_1:config",
			},
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-auth-msvc/health"},
		Component{
			name: "case-flow",
			image: "orchard/orchard-case-flow-msvc:latest",
			dockerId: "dcmp_orchard-case-flow-msvc_1",
			links: []string {
				"dcmp_orchard-config-msvc_1:redis",
				"dcmp_orchard-doc-analysis-msvc_1:config",
			},
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-case-flow-msvc/health"},
		Component{
			name: "doc-analysis",
			image: "orchard/orchard-doc-analysis-msvc:latest",
			dockerId: "dcmp_orchard-doc-analysis-msvc_1",
			links: []string {
				"dcmp_orchard-config-msvc_1:config",
			},
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-doc-analysis-msvc/health"},
		Component{
			name: "ui",
			image: "orchard/orchard-doc-analysis-ui:latest",
			dockerId: "dcmp_orchard-doc-analysis-ui_1",
			containerPort: 80,
			hostPort: 3000,
			testUrl:  "http://localhost:3000/"},
	}
}
