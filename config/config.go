package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Logger struct {
	Debug bool `yaml:"debug"`
}

type HTTP struct {
	Port            string `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	WriteTimeout    int    `yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
	ReadTimeout     int    `yaml:"readTimeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
	ShutdownTimeout int    `yaml:"shutdownTimeout" env:"HTTP_SHUT_DOWN_TIMEOUT" env-required:"true"`
}

type Postgres struct {
	User         string `env:"POSTGRES_USER"`
	Password     string `env:"POSTGRES_PASSWORD"`
	Host         string `env:"POSTGRES_HOST"`
	Port         string `env:"POSTGRES_PORT"`
	DB           string `env:"POSTGRES_DB"`
	SSLMode      string `env:"POSTGRES_SSLMODE"`
	ConnAttempts int    `yaml:"connAttempts"`
	ConnTimeout  int    `yaml:"connTimeout"`
}

type Config struct {
	Logger   `yaml:"logger"`
	HTTP     `yaml:"http"`
	Postgres `yaml:"postgres"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config/development.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
