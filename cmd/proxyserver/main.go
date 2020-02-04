package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/sv-z/amqproxy/Internal/app/config"
	"github.com/sv-z/amqproxy/Internal/app/proxyserver"
	"os"
)

func main() {
	conf, err := config.NewConfig()

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	if err := proxyserver.Start(conf); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
