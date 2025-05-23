package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CsKey      string `envconfig:"CLOUDSQUID_API_KEY" required:"true"`
	CsEndpoint string `envconfig:"CLOUDSQUID_API_ENDPOINT" required:"true"`
	CsSourceID string `envconfig:"CLOUDSQUID_AGENT_ID" required:"true"`
}

func Load() (*Config, error) {
	if err := godotenv.Overload(); err != nil {
		return nil, fmt.Errorf("loading environmental variables: %w", err)
	}

	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("processing env config %w", err)
	}

	log.Printf("Loaded config: %s", config)

	return &config, nil
}
