// Package config
package config

import (
	"log"
	rest "vk-inter/internal/transport"
	"vk-inter/pkg/db/mongo"

	"github.com/ilyakaznacheev/cleanenv"
)

type KeyString string

const ConfigKey KeyString = "config"

// Config
type Config struct {
	mongo.MongoConfig
	rest.RestConfig
	Debug  bool   `env:"DEBUG" env-default:"true"`
	Secret string `env:"SECRET" env-default:"test_key"`
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
