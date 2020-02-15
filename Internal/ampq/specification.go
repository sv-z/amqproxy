package ampq

import (
	"bufio"
	"fmt"
	"io"
)

type Specification interface {
	// The check protocol version
	check(hello []byte) bool

	// Get request data to start connect
	SendConnectionStart() bool
	ReceiveConnectionStartOk() bool
	SendConnectionTune() bool
	ReceiveConnectionTuneOK() bool
	ReceiveConnectionOpen() bool
	SendConnectionOpenOK() bool
}

func GetSpecification(rw io.ReadWriter) (error, Specification) {

	protocolHeader := make([]byte, 8)
	reader := bufio.NewReader(rw)
	if _, err := reader.Read(protocolHeader); err == io.EOF {
		return err, nil
	}

	if (Specification91{}).check(protocolHeader) {
		return nil, newSpecification91(rw)
	}

	return fmt.Errorf("unsupported protocol: '%v'", protocolHeader), nil

}
