package ampq

import (
	"bufio"
	"bytes"
	"fmt"
	transfer "github.com/sv-z/amqproxy/Internal/ampq/data-transfer"
	"github.com/sv-z/amqproxy/Internal/ampq/spec091"
	"io"
)

type Connection struct {
	Connected bool
	rw        *io.ReadWriter
	spec      *spec091.Spec
}

func NewConnection(readWriter io.ReadWriter) *Connection {
	return &Connection{
		rw:   &readWriter,
		spec: spec091.NewSpec091(&readWriter),
	}
}

//    Connection          = open-Connection *use-Connection close-Connection
//    open-Connection     = C:protocol-header
//                          S:START C:START-OK
//                          *challenge
//                          S:TUNE C:TUNE-OK
//                          C:OPEN S:OPEN-OK
//    challenge           = S:SECURE C:SECURE-OK
//    use-Connection      = *channel
//    close-Connection    = C:CLOSE S:CLOSE-OK
//                        / S:CLOSE C:CLOSE-OK
func (c Connection) Open() error {

	//The client MUST start a new connection by sending a protocol header.
	// C:protocol-header
	protocolHeader := make([]byte, 8)
	reader := bufio.NewReader(*c.rw)
	if _, err := reader.Read(protocolHeader); err == io.EOF {
		return err
	}

	if !c.checkProtocol(protocolHeader) {
		return fmt.Errorf("unsupported protocol: '%v'", protocolHeader)
	}

	// S:START >> C:START-OK
	if err := c.startOk(); err != nil {
		return err
	}

	// S:TUNE C:TUNE-OK
	if err := c.tuneOK(); err != nil {
		return err
	}

	// C:OPEN S:OPEN-OK
	if err := c.openOK(); err != nil {
		return err
	}

	c.Connected = true

	return nil
}

// Server send connection.start
// The client return connection.start-ok
func (c *Connection) startOk() error {
	if !c.sendStart() {
		return fmt.Errorf("cannot send response \"connection.start\"")
	}

	if _, err := c.spec.PullConnectionStartOk(); err != nil {
		return err
	}

	return nil
}

// The server send connection.tune
// The client return connection.tune-ok
func (c *Connection) tuneOK() error {

	if !c.spec.PushConnectionTune() {
		return fmt.Errorf("cannot send response \"connection.tune-ok\"")
	}

	if _, err := c.spec.PullConnectionTuneOK(); err != nil {
		return err
	}

	return nil
}

// The client send connection open
// The server return connection open-ok
func (c *Connection) openOK() error {

	if _, err := c.spec.PullConnectionOpen(); err != nil {
		panic(err)
	}

	if !c.spec.PushConnectionOpenOK() {
		return fmt.Errorf("cannot send response \"connection.open-ok\"")
	}

	return nil
}

func (s Connection) checkProtocol(protocolHeader []byte) bool {
	return bytes.Compare(protocolHeader, []byte{'A', 'M', 'Q', 'P', 0, 0, 9, 1}) == 0
}

func (c *Connection) sendStart() bool {

	args := transfer.Table{
		"capabilities": transfer.Table{
			"consumer_priorities":          true,
			"authentication_failure_close": true,
			"direct_reply_to":              true,
			"publisher_confirms":           true,
			"exchange_exchange_bindings":   true,
			"basic.nack":                   true,
			"consumer_cancel_notify":       true,
			"connection.blocked":           true,
			"per_consumer_qos":             true,
		},
		"cluster_name": "rabbit@rabbit",
		"copyright":    "Copyright (c) 2007-2019 Pivotal Software, Inc.",
		"information":  "Licensed under the MPL 1.1. Website: https://rabbitmq.com",
		"platform":     "Erlang/OTP 22.2.4",
		"product":      "RabbitMQ",
		"version":      "3.8.2",
	}

	return c.spec.PushConnectionStart(&args)
}
