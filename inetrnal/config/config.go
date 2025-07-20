// Package config
package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config
type Config struct {
	Debug bool `env:"DEBUG" env-default:"true"`
}

// New
func New() *Config {
	cfg := Config{}

	err := cleanenv.ReadConfig("./.env", &cfg)

	if err != nil {
		log.Fatalf("error reading config: %v`", err)
	}

	return &cfg
}
