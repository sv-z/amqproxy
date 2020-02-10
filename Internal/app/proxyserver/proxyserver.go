package proxyserver

import (
	"fmt"
	guuid "github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"github.com/sv-z/amqproxy/Internal/ampq"
	"github.com/sv-z/amqproxy/Internal/app/config"
	"net"
)

// Start server
func Start(conf *config.Config) *error {
	setLoggerLevel(conf)

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

// Set mail logger level
func setLoggerLevel(conf *config.Config) {
	level, err := logger.ParseLevel(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	logger.SetLevel(level)
	logger.WithField("requestId", "1")
}

// Handle request
func handleRequest(conn net.Conn) {
	requestId := guuid.New().String()

	logger.Debug(fmt.Sprintf("----- ==== Start connection [%s] ==== -----", requestId))
	defer func() {
		logger.Debug(fmt.Sprintf("----- ==== Stop connection [%s] ==== -----", requestId))
	}()

	err, sp := ampq.GetSpecification(conn)
	if err != nil {
		logger.Error(fmt.Sprintf("The %s", err))
		return
	}

	if !sp.SendResponseConnectionStart() {
		logger.Error("Cannot send response \"connection.start\".")
		return
	}
	logger.Debug(fmt.Sprintf("----- ==== send connection.start [%s] ==== -----", requestId))

	if !sp.WaitConnectionStartOk() {
		logger.Error("Clients request \"connection.start­ok\" not received.")
	}

	// ~/proj/home/go-lang/server
}
