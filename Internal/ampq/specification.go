package ampq

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Specification interface {
	// The check protocol version
	check(hello []byte) bool

	// Get request data to start connect
	SendResponseConnectionStart() bool
	WaitConnectionStartOk() bool
}

func GetSpecification(rw io.ReadWriter) (error, Specification) {

	hello := make([]byte, 8)
	reader := bufio.NewReader(rw)
	if _, err := reader.Read(hello); err == io.EOF {
		return err, nil
	}

	if (Specification91{}).check(hello) {
		return nil, newSpecification91(rw)
	}

	return fmt.Errorf("unsupported protocol: '%v'", hello), nil

}

// Prepare method header
func prepareMethod(args ...interface{}) []byte {
	var payload bytes.Buffer

	for _, val := range args {
		if err := binary.Write(&payload, binary.BigEndian, val); err != nil {
			return nil
		}
	}

	return payload.Bytes()
}

//
//func writeFrame(w io.Writer, typ uint8, channel uint16, payload []byte) (err error) {
//	end := []byte{frameEnd}
//	size := uint(len(payload))
//
//	_, err = w.Write([]byte{
//		byte(typ),
//		byte((channel & 0xff00) >> 8),
//		byte((channel & 0x00ff) >> 0),
//		byte((size & 0xff000000) >> 24),
//		byte((size & 0x00ff0000) >> 16),
//		byte((size & 0x0000ff00) >> 8),
//		byte((size & 0x000000ff) >> 0),
//	})
//
//	if err != nil {
//		return
//	}
//
//	if _, err = w.Write(payload); err != nil {
//		return
//	}
//
//	if _, err = w.Write(end); err != nil {
//		return
//	}
//
//	return
//}
