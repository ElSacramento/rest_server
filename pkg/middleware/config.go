package middleware

import (
	"io/ioutil"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port" validate:"required"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Server   ServerConfig   `yaml:"server" validate:"required"`
}

func ParseConfig(cfgPath *string) (*Config, error) {
	data, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) ValidateConfig() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}
