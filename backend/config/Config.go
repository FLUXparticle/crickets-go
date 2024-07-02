package config

import (
	"os"
	"strings"
)

type Config struct {
	Hostname string
	ApiKey   string
}

func NewConfig() *Config {
	hostname, _ := os.Hostname()
	if strings.Contains(hostname, ".") {
		hostname = "localhost"
	}

	apiKey := os.Getenv("API_KEY")

	return &Config{
		Hostname: hostname,
		ApiKey:   apiKey,
	}
}
