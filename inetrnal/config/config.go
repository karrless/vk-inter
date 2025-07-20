// Package config
package config

import (
	"log"
	"vk-inter/pkg/db/mongo"

	"github.com/ilyakaznacheev/cleanenv"
)

type KeyString string

const ConfigKey KeyString = "config"

// Config
type Config struct {
	mongo.MongoConfig
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
