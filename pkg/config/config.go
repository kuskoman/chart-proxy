package config

type Config struct {
	Server   *ServerConfig      `hcl:"server,block"`
	Logging  *LoggingConfig     `hcl:"logging,block"`
	Aliases  *[]RepositoryAlias `hcl:"repository_alias,block"`
	Mappings *[]Mapping         `hcl:"mapping,block"`
}

type ServerConfig struct {
	Port int    `hcl:"port"`
	Host string `hcl:"host"`
}

type LoggingConfig struct {
	Level  string `hcl:"level"`
	Format string `hcl:"format"`
}

type RepositoryAlias struct {
	Name string `hcl:"name,label"`
	URL  string `hcl:"url"`
}

type Mapping struct {
	Name       string `hcl:"name,label"`
	Upstream   string `hcl:"upstream"`
	Downstream string `hcl:"downstream"`
}
