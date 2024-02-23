package config

func getDefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Logging: &LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Aliases:  &[]RepositoryAlias{},
		Mappings: &[]Mapping{},
	}
}
