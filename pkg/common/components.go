package common

var components []Component

type Component struct {
	Name          string   `yaml:"name,omitempty"`
	DockerId      string   `yaml:"dockerId,omitempty"`
	TestUrl       string   `yaml:"testUrl,omitempty"`
	Image         string   `yaml:"image,omitempty"`
	ContainerPort int      `yaml:"containerPort,omitempty"`
	HostPort      int      `yaml:"hostPort,omitempty"`
	Repository    string   `yaml:"repository,omitempty"`
	Env           []string `yaml:"env,omitempty"`
	Links         []string `yaml:"links,omitempty"`
	Mounts        []string `yaml:"mounts,omitempty"`
}

func ComponentNames(components []Component) []string {
	componentNames := []string{}
	for _, component := range components {
		componentNames = append(componentNames, component.Name)
	}
	return componentNames
}

func ComponentMap(components []Component) map[string]Component {
	elementMap := make(map[string]Component)
	for i := 0; i < len(components); i += 1 {
		elementMap[components[i].Name] = components[i]
	}
	return elementMap
}
