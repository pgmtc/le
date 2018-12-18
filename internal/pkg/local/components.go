package local

type Component struct {
	name     string
	dockerId string
	testUrl  string
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
			name:     "db",
			dockerId: "dcmp_orchard-local-db_1",
			testUrl:  ""},
		Component{
			name:     "redis",
			dockerId: "dcmp_orchard-redis_1",
			testUrl:  ""},
		Component{
			name:     "config",
			dockerId: "dcmp_orchard-config-msvc_1",
			testUrl:  ""},
		Component{
			name:     "gateway",
			dockerId: "dcmp_orchard-gateway-msvc_1",
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/health"},
		Component{
			name:     "auth",
			dockerId: "dcmp_orchard-auth-msvc_1",
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-auth-msvc/health"},
		Component{
			name:     "case-flow",
			dockerId: "dcmp_orchard-case-flow-msvc_1",
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-case-flow-msvc/health"},
		Component{
			name:     "doc-analysis-msvc",
			dockerId: "dcmp_orchard-doc-analysis-msvc_1",
			testUrl:  "http://localhost:8765/orchard-gateway-msvc/orchard-doc-analysis-msvc/health"},
		Component{
			name:     "doc-analysis-ui",
			dockerId: "dcmp_orchard-doc-analysis-ui_1",
			testUrl:  "http://localhost:3000/"},
	}
}
