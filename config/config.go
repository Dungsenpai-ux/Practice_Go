package config

type Config struct {
	Version string
}

func NewConfig() *Config {
	return &Config{
		Version: "v1.0.0",
	}
}