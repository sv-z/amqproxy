package ampq

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	frameMethod        = 1
	frameHeader        = 2
	frameBody          = 3
	frameHeartbeat     = 8
	frameMinSize       = 4096
	frameEnd           = 206 // "\xCE"
	replySuccess       = 200
	ContentTooLarge    = 311
	NoRoute            = 312
	NoConsumers        = 313
	ConnectionForced   = 320
	InvalidPath        = 402
	AccessRefused      = 403
	NotFound           = 404
	ResourceLocked     = 405
	PreconditionFailed = 406
	FrameError         = 501
	SyntaxError        = `502`
	CommandInvalid     = 503
	ChannelError       = 504
	UnexpectedFrame    = 505
	ResourceError      = 506
	NotAllowed         = 530
	NotImplemented     = 540
	InternalError      = 541
)

func newSpecification91(readWriter io.ReadWriter) *Specification91 {
	sp := &Specification91{
		readWriter:   readWriter,
		versionMajor: byte(0),
		versionMinor: byte(9),
		locales:      "en_US",
		mechanisms:   "PLAIN",
	}

	return sp
}

type Specification91 struct {
	readWriter   io.ReadWriter
	versionMajor byte
	versionMinor byte
	locales      string
	mechanisms   string
}

func (s Specification91) check(hello []byte) bool {
	return bytes.Compare(hello, []byte{'A', 'M', 'Q', 'P', 0, 0, 9, 1}) == 0
}

func (sp *Specification91) SendResponseConnectionStart() bool {

	args := Table{
		"capabilities": Table{
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

	payload := prepareMethod(
		uint16(10), //class,
		uint16(10), //method,
		sp.versionMajor,
		sp.versionMinor,
		mapToByte(args),
		longStrToByte(sp.mechanisms),
		longStrToByte(sp.locales),
	)
	if payload == nil {
		return false
	}

	return writeFrame(sp.readWriter, frameMethod, 0, payload) == nil
}

func (sp *Specification91) WaitConnectionStartOk() bool {
	return false
}

func longStrToByte(str string) []byte {
	b := []byte(str)
	var length = uint32(len(b))
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, length); err != nil {
		panic(err)
	}

	if _, err := buf.Write(b[:length]); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func mapToByte(table Table) []byte {
	var buf bytes.Buffer

	for key, val := range table {
		if err := writeShortstr(&buf, key); err != nil {
			panic(err)
		}
		if err := writeField(&buf, val); err != nil {
			panic(err)
		}
	}

	return longStrToByte(string(buf.Bytes()))
}
