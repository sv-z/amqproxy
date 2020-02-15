package spec091

import transfer "github.com/sv-z/amqproxy/Internal/ampq/data-transfer"

type methodFrame struct {
	ClassId  uint16
	MethodId uint16
	Method   interface{}
}

// ------------------------------------------ CLASS 10 -----------------------------------------------------------------
type connectionStartOk struct {
	ClientProperties transfer.Table
	Mechanism        string
	Response         string
	Locale           string
}

type connectionTuneOk struct {
	ChannelMax uint16
	FrameMax   uint32
	Heartbeat  uint16
}

type connectionOpen struct {
	VirtualHost string
	reserved1   string
	reserved2   bool
}
