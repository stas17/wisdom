package config

// Config holds all the configuration data.
type Config struct {
	App     *App     `yaml:"app"`
	Logging *Logging `yaml:"logging"`
}

// App holds the application configuration.
type App struct {
	HTTPWisdomAddr string `yaml:"httpWisdomAddr"`
	Salt           string `yaml:"salt"`
}

// Logging holds the logging configuration.
type Logging struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`
}
