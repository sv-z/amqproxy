package config

import (
	"fmt"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found")
	}
}

type Config struct {
	BindAddr string
	BindPort int
	LogLevel string
}

// Create new app config
func NewConfig() (*Config, error) {
	proxyHost, exists := os.LookupEnv("PROXY_CONNECTION_HOST")
	if !exists {
		return nil, fmt.Errorf(`parameter "PROXY_CONNECTION_HOST" not found in .env file`)
	}
	proxyPortDraft, exists := os.LookupEnv("PROXY_CONNECTION_PORT")
	if !exists {
		return nil, fmt.Errorf(`parameter "PROXY_CONNECTION_PORT" not found in .env file`)
	}

	proxyPort, err := strconv.Atoi(proxyPortDraft)
	if err != nil {
		return nil, fmt.Errorf(`parameter "PROXY_CONNECTION_PORT" must be integer`)
	}

	logLevel, exists := os.LookupEnv("LOG_LEVEL")
	if !exists {
		logLevel = "error"
	}

	return &Config{
		BindAddr: proxyHost,
		BindPort: proxyPort,
		LogLevel: logLevel,
	}, nil
}
