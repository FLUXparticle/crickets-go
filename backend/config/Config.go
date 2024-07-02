package config

import (
	"os"
	"strings"
)

type Config struct {
	Hostname string
	ApiKey   string
	AmqpHost string
}

func NewConfig() *Config {
	hostname, _ := os.Hostname()
	if strings.Contains(hostname, ".") {
		hostname = ""
	}

	apiKey := os.Getenv("API_KEY")

	amqpHost := os.Getenv("AMQP_HOST")
	if len(amqpHost) == 0 {
		amqpHost = "localhost"
	}

	return &Config{
		Hostname: hostname,
		ApiKey:   apiKey,
		AmqpHost: amqpHost,
	}
}
