package data_transfer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type dataReader struct {
	r *io.Reader
}

func NewDataReader(r io.Reader) *dataReader {
	return &dataReader{r: &r}
}

func (r *dataReader) ReadTable() (table Table, err error) {
	return readTable(*r.r)
}

func (r *dataReader) ReadShortstr() (v string, err error) {
	return readShortstr(*r.r)
}

func (r *dataReader) ReadLongstr() (v string, err error) {
	return readLongstr(*r.r)
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
func readField(r io.Reader) (v interface{}, err error) {
	var typ byte
	if err = binary.Read(r, binary.BigEndian, &typ); err != nil {
		return
	}

	switch typ {
	case 't':
		var value uint8
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return (value != 0), nil

	case 'b':
		var value [1]byte
		if _, err = io.ReadFull(r, value[0:1]); err != nil {
			return
		}
		return value[0], nil

	case 's':
		var value int16
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return value, nil

	case 'I':
		var value int32
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return value, nil

	case 'l':
		var value int64
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return value, nil

	case 'f':
		var value float32
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return value, nil

	case 'd':
		var value float64
		if err = binary.Read(r, binary.BigEndian, &value); err != nil {
			return
		}
		return value, nil

	case 'D':
		return readDecimal(r)

	case 'S':
		return readLongstr(r)

	case 'A':
		return readArray(r)

	case 'T':
		return readTimestamp(r)

	case 'F':
		return readTable(r)

	case 'x':
		var len int32
		if err = binary.Read(r, binary.BigEndian, &len); err != nil {
			return nil, err
		}

		value := make([]byte, len)
		if _, err = io.ReadFull(r, value); err != nil {
			return nil, err
		}
		return value, err

	case 'V':
		return nil, nil
	}

	return nil, fmt.Errorf("invalid field or value inside of a frame")
}

/*
	Field tables are long strings that contain packed name-value pairs.  The
	name-value pairs are encoded as short string defining the name, and octet
	defining the values type and then the value itself.   The valid field types for
	tables are an extension of the native integer, bit, string, and timestamp
	types, and are shown in the grammar.  Multi-octet integer fields are always
	held in network byte order.
*/
func readTable(r io.Reader) (table Table, err error) {
	var nested bytes.Buffer
	var str string

	if str, err = readLongstr(r); err != nil {
		return
	}

	nested.Write([]byte(str))

	table = make(Table)

	for nested.Len() > 0 {
		var key string
		var value interface{}

		if key, err = readShortstr(&nested); err != nil {
			return
		}

		if value, err = readField(&nested); err != nil {
			return
		}

		table[key] = value
	}

	return
}

func readArray(r io.Reader) ([]interface{}, error) {
	var (
		size uint32
		err  error
	)

	if err = binary.Read(r, binary.BigEndian, &size); err != nil {
		return nil, err
	}

	var (
		lim   = &io.LimitedReader{R: r, N: int64(size)}
		arr   = []interface{}{}
		field interface{}
	)

	for {
		if field, err = readField(lim); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		arr = append(arr, field)
	}

	return arr, nil
}

func readShortstr(r io.Reader) (v string, err error) {
	var length uint8
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}

	bytes := make([]byte, length)
	if _, err = io.ReadFull(r, bytes); err != nil {
		return
	}
	return string(bytes), nil
}

func readLongstr(r io.Reader) (v string, err error) {
	var length uint32
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}

	// slices can't be longer than max int32 value
	if length > (^uint32(0) >> 1) {
		return
	}

	bytes := make([]byte, length)
	if _, err = io.ReadFull(r, bytes); err != nil {
		return
	}
	return string(bytes), nil
}

func readDecimal(r io.Reader) (v Decimal, err error) {
	if err = binary.Read(r, binary.BigEndian, &v.Scale); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &v.Value); err != nil {
		return
	}
	return
}

func readFloat32(r io.Reader) (v float32, err error) {
	if err = binary.Read(r, binary.BigEndian, &v); err != nil {
		return
	}
	return
}

func readFloat64(r io.Reader) (v float64, err error) {
	if err = binary.Read(r, binary.BigEndian, &v); err != nil {
		return
	}
	return
}

func readTimestamp(r io.Reader) (v time.Time, err error) {
	var sec int64
	if err = binary.Read(r, binary.BigEndian, &sec); err != nil {
		return
	}
	return time.Unix(sec, 0), nil
}
