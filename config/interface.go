package config

type Repository interface {
	LoadFile(fileName string) error
	SaveFile(fileName string) error
	GetConfig() *Config
}

type Service interface {
	LoadConfig(fileName string) (*Config, error)
}
