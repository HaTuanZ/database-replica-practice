package compose

type ComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Image         string            `yaml:"image"`
	Command       string            `yaml:"command"`
	ContainerName string            `yaml:"container_name"`
	Environment   map[string]string `yaml:"environment"`
	Ports         []string          `yaml:"ports"`
	Volumes       []string          `yaml:"volumes"`
}
