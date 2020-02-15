package spec091

import (
	"fmt"
	transfer "github.com/sv-z/amqproxy/Internal/ampq/data-transfer"
	"io"
)

func NewSpec091(readWriter io.ReadWriter) *Spec {
	return &Spec{
		readWriter:   readWriter,
		versionMajor: byte(0),
		versionMinor: byte(9),
		locales:      "en_US",
		mechanisms:   "PLAIN",
	}
}

type Spec struct {
	readWriter   io.ReadWriter
	versionMajor byte
	versionMinor byte
	locales      string
	mechanisms   string
}

func (spec *Spec) PushConnectionStart(args *transfer.Table) bool {
	payload := prepareMethod(
		uint16(10), //class,
		uint16(10), //method,
		spec.versionMajor,
		spec.versionMinor,
		transfer.MapToByte(*args),
		transfer.LongStrToByte(spec.mechanisms),
		transfer.LongStrToByte(spec.locales),
	)
	if payload == nil {
		return false
	}

	return writeFrame(spec.readWriter, frameMethod, 0, payload) == nil
}

func (spec *Spec) PullConnectionStartOk() (*connectionStartOk, error) {

	mfi, err := readFrame(spec.readWriter)
	if err != nil {
		return nil, err
	}

	if mf, ok := mfi.(*methodFrame); ok {
		if resp, ok := mf.Method.(*connectionStartOk); ok {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("invalid C:START-OK receive")
}

// connection.tune
func (spec *Spec) PushConnectionTune() bool {
	payload := prepareMethod(
		uint16(10),     //class,
		uint16(30),     //method
		uint16(2047),   //ChannelMax
		uint32(131072), //FrameMax
		uint16(60),     //Heartbeat
	)
	if payload == nil {
		return false
	}

	return writeFrame(spec.readWriter, frameMethod, 0, payload) == nil
}

func (spec *Spec) PullConnectionTuneOK() (*connectionTuneOk, error) {

	mfi, err := readFrame(spec.readWriter)
	if err != nil {
		return nil, err
	}

	if mf, ok := mfi.(*methodFrame); ok {
		if resp, ok := mf.Method.(*connectionTuneOk); ok {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("invalid C:TUNE-OK receive")
}

func (spec *Spec) PullConnectionOpen() (*connectionOpen, error) {
	mfi, err := readFrame(spec.readWriter)
	if err != nil {
		return nil, err
	}

	if mf, ok := mfi.(*methodFrame); ok {
		if resp, ok := mf.Method.(*connectionOpen); ok {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("invalid C:OPEN receive")
}

func (spec *Spec) PushConnectionOpenOK() bool {
	payload := prepareMethod(
		uint16(10),                  //class,
		uint16(41),                  //method
		transfer.ShortStrToByte(""), // reserved1
	)
	if payload == nil {
		return false
	}

	return writeFrame(spec.readWriter, frameMethod, 0, payload) == nil
}
