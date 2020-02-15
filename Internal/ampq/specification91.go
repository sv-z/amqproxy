package ampq

import (
	"bytes"
	transfer "github.com/sv-z/amqproxy/Internal/ampq/data-transfer"
	"github.com/sv-z/amqproxy/Internal/ampq/spec091"
	"io"
)

func newSpecification91(readWriter io.ReadWriter) *Specification91 {
	sp := &Specification91{
		spec091.NewSpec091(readWriter),
	}

	return sp
}

type Specification91 struct {
	*spec091.Spec
}

func (s Specification91) check(hello []byte) bool {
	return bytes.Compare(hello, []byte{'A', 'M', 'Q', 'P', 0, 0, 9, 1}) == 0
}

func (sp *Specification91) SendConnectionStart() bool {

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

	return sp.PushConnectionStart(&args)
}

func (sp *Specification91) ReceiveConnectionStartOk() bool {
	if _, err := sp.PullConnectionStartOk(); err != nil {
		panic(err)
	}

	return true
}

func (sp *Specification91) SendConnectionTune() bool {
	return sp.PushConnectionTune()
}

func (sp *Specification91) ReceiveConnectionTuneOK() bool {
	if _, err := sp.PullConnectionTuneOK(); err != nil {
		panic(err)
	}

	return true
}

func (sp *Specification91) ReceiveConnectionOpen() bool {
	if _, err := sp.PullConnectionOpen(); err != nil {
		panic(err)
	}

	return true
}

func (sp *Specification91) SendConnectionOpenOK() bool {

	return sp.PushConnectionOpenOK()
}
