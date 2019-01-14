package common

var components []Component

type Component struct {
	Name          string   `yaml:"name,omitempty"`
	DockerId      string   `yaml:"dockerId,omitempty"`
	TestUrl       string   `yaml:"testUrl,omitempty"`
	Image         string   `yaml:"image,omitempty"`
	ContainerPort int      `yaml:"containerPort,omitempty"`
	HostPort      int      `yaml:"hostPort,omitempty"`
	Env           []string `yaml:"env,omitempty"`
	Links         []string `yaml:"links,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	DockerFile    string   `yaml:"dockerFile,omitempty"`
	BuildRoot     string   `yaml:"buildRoot,omitempty"`
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
	config := FileSystemConfig("~/.orchard-zzz")
	return config.CurrentProfile().Components
}
