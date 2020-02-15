package data_transfer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

func LongStrToByte(str string) []byte {
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

func ShortStrToByte(str string) []byte {
	b := []byte(str)
	var length = uint8(len(b))
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, length); err != nil {
		panic(err)
	}

	if _, err := buf.Write(b[:length]); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func MapToByte(table Table) []byte {
	var buf bytes.Buffer

	for key, val := range table {
		if err := writeShortstr(&buf, key); err != nil {
			panic(err)
		}
		if err := writeField(&buf, val); err != nil {
			panic(err)
		}
	}

	return LongStrToByte(string(buf.Bytes()))
}

// --------------------------------------------------------------------------

func writeShortstr(w io.Writer, s string) (err error) {
	b := []byte(s)

	var length = uint8(len(b))

	if err = binary.Write(w, binary.BigEndian, length); err != nil {
		return
	}

	if _, err = w.Write(b[:length]); err != nil {
		return
	}

	return
}

func writeLongstr(w io.Writer, s string) (err error) {
	b := []byte(s)

	var length = uint32(len(b))

	if err = binary.Write(w, binary.BigEndian, length); err != nil {
		return
	}

	if _, err = w.Write(b[:length]); err != nil {
		return
	}

	return
}

/*
'A': []interface{}
'D': Decimal
'F': Table
'I': int32
'S': string
'T': time.Time
'V': nil
'b': byte
'd': float64
'f': float32
'l': int64
's': int16
't': bool
'x': []byte
*/
func writeField(w io.Writer, value interface{}) (err error) {
	var buf [9]byte
	var enc []byte

	switch v := value.(type) {
	case bool:
		buf[0] = 't'
		if v {
			buf[1] = byte(1)
		} else {
			buf[1] = byte(0)
		}
		enc = buf[:2]

	case byte:
		buf[0] = 'b'
		buf[1] = byte(v)
		enc = buf[:2]

	case int16:
		buf[0] = 's'
		binary.BigEndian.PutUint16(buf[1:3], uint16(v))
		enc = buf[:3]

	case int:
		buf[0] = 'I'
		binary.BigEndian.PutUint32(buf[1:5], uint32(v))
		enc = buf[:5]

	case int32:
		buf[0] = 'I'
		binary.BigEndian.PutUint32(buf[1:5], uint32(v))
		enc = buf[:5]

	case int64:
		buf[0] = 'l'
		binary.BigEndian.PutUint64(buf[1:9], uint64(v))
		enc = buf[:9]

	case float32:
		buf[0] = 'f'
		binary.BigEndian.PutUint32(buf[1:5], math.Float32bits(v))
		enc = buf[:5]

	case float64:
		buf[0] = 'd'
		binary.BigEndian.PutUint64(buf[1:9], math.Float64bits(v))
		enc = buf[:9]

	case Decimal:
		buf[0] = 'D'
		buf[1] = byte(v.Scale)
		binary.BigEndian.PutUint32(buf[2:6], uint32(v.Value))
		enc = buf[:6]

	case string:
		buf[0] = 'S'
		binary.BigEndian.PutUint32(buf[1:5], uint32(len(v)))
		enc = append(buf[:5], []byte(v)...)

	case []interface{}: // field-array
		buf[0] = 'A'

		sec := new(bytes.Buffer)
		for _, val := range v {
			if err = writeField(sec, val); err != nil {
				return
			}
		}

		binary.BigEndian.PutUint32(buf[1:5], uint32(sec.Len()))
		if _, err = w.Write(buf[:5]); err != nil {
			return
		}

		if _, err = w.Write(sec.Bytes()); err != nil {
			return
		}

		return

	case time.Time:
		buf[0] = 'T'
		binary.BigEndian.PutUint64(buf[1:9], uint64(v.Unix()))
		enc = buf[:9]

	case Table:
		if _, err = w.Write([]byte{'F'}); err != nil {
			return
		}
		return writeTable(w, v)

	case []byte:
		buf[0] = 'x'
		binary.BigEndian.PutUint32(buf[1:5], uint32(len(v)))
		if _, err = w.Write(buf[0:5]); err != nil {
			return
		}
		if _, err = w.Write(v); err != nil {
			return
		}
		return

	case nil:
		buf[0] = 'V'
		enc = buf[:1]

	default:
		return fmt.Errorf("unsupported table field type")
	}

	_, err = w.Write(enc)

	return
}

func writeTable(w io.Writer, table Table) (err error) {
	var buf bytes.Buffer

	for key, val := range table {
		if err = writeShortstr(&buf, key); err != nil {
			return
		}
		if err = writeField(&buf, val); err != nil {
			return
		}
	}

	return writeLongstr(w, string(buf.Bytes()))
}
