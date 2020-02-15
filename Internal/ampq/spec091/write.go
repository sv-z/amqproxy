package spec091

import (
	"bytes"
	"encoding/binary"
	"io"
)

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

func writeFrame(w io.Writer, typ uint8, channel uint16, payload []byte) (err error) {
	end := []byte{frameEnd}
	size := uint(len(payload))

	_, err = w.Write([]byte{
		byte(typ),
		byte((channel & 0xff00) >> 8),
		byte((channel & 0x00ff) >> 0),
		byte((size & 0xff000000) >> 24),
		byte((size & 0x00ff0000) >> 16),
		byte((size & 0x0000ff00) >> 8),
		byte((size & 0x000000ff) >> 0),
	})

	if err != nil {
		return
	}

	if _, err = w.Write(payload); err != nil {
		return
	}

	if _, err = w.Write(end); err != nil {
		return
	}

	return
}
