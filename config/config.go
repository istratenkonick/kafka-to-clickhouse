package config

type KafkaConfig struct {
	Host []string `yaml:"host"`
}

type ClickhouseConfig struct {
	Database       string `yaml:"database"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Authentication bool   `yaml:"auth"`
	Host           string `yaml:"host"`
}

type CoreConfig struct {
	Files map[string]File `yaml:"files"`
	Host  string          `yaml:"host"`
}

type Config struct {
	Core       *CoreConfig
	Clickhouse *ClickhouseConfig
}

type File struct {
	Enabled bool     `yaml:"enabled"`
	Topics  []string `yaml:"topics"`
}
