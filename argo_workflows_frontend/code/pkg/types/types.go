package types

type DomainMapping struct {
	Production  map[string]string `yaml:"main"`
	Development map[string]string `yaml:"development"`
}

type EnvironmentMapping map[string]EnvironmentConfig

type EnvironmentConfig struct {
	Bucket       string `yaml:"bucket"`
	Distribution string `yaml:"distribution"`
	Domain       string `yaml:"domain"`
}

type ProjectConfig struct {
	ReleaseStage string            `yaml:"releaseStage"`
	Urls         map[string]string `yaml:"urls"`
}

type Stages map[string]map[string]map[string]string
