package proxyserver

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/sv-z/amqproxy/Internal/app/config"
	"net"
)

// Start server
func Start(conf *config.Config) *error {

	address := fmt.Sprintf("%s:%d", conf.BindAddr, conf.BindPort)
	// Listen for incoming connections.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return &err
	}

	// Close the listener when the application closes.
	defer listener.Close()
	logger.Info(fmt.Sprintf(`Listening on tcp: "%s"`, address))

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Warn(fmt.Sprintf(`Error accepting: %s`, err.Error()))
			continue
		}

		go handleRequest(conn)
	}
}

// Handle request
func handleRequest(conn net.Conn) {
	logger.Warn(fmt.Sprintf(`----------****--------------------`))
	// ~/proj/home/go-lang/server
}
