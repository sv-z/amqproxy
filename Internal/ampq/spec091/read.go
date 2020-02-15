package spec091

import (
	"encoding/binary"
	"fmt"
	transfer "github.com/sv-z/amqproxy/Internal/ampq/data-transfer"
	"io"
)

func readFrame(reader io.Reader) (fm interface{}, err error) {
	var scratch [7]byte

	if _, err = io.ReadFull(reader, scratch[:7]); err != nil {
		return
	}

	typ := uint8(scratch[0])
	channel := binary.BigEndian.Uint16(scratch[1:3])
	size := binary.BigEndian.Uint32(scratch[3:7])

	fmt.Println(typ, "  ", channel, "   ", size)

	switch typ {
	case frameMethod:
		if fm, err = parseMethodFrame(reader); err != nil {
			return
		}

	case frameHeader:
		//if frame, err = r.parseHeaderFrame(channel, size); err != nil {
		//	return
		//}

	case frameBody:
		//if frame, err = r.parseBodyFrame(channel, size); err != nil {
		//	return nil, err
		//}

	case frameHeartbeat:
		//if frame, err = r.parseHeartbeatFrame(channel, size); err != nil {
		//	return
		//}

	default:
		return fm, fmt.Errorf("frame could not be parsed")
	}

	if _, err = io.ReadFull(reader, scratch[:1]); err != nil {
		return fm, err
	}

	if scratch[0] != frameEnd {
		return fm, fmt.Errorf("frame could not be parsed")
	}

	return
}

func parseMethodFrame(reader io.Reader) (mf *methodFrame, err error) {

	mf = &methodFrame{}

	dataReader := transfer.NewDataReader(reader)

	var classId uint16
	var methodId uint16

	if err = binary.Read(reader, binary.BigEndian, &classId); err != nil {
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &methodId); err != nil {
		return
	}

	switch classId {
	case 10: // connection
		switch methodId {
		//case 10: // connection start
		case 11: // connection start-ok
			method := connectionStartOk{}
			mf.Method = &method
			if method.ClientProperties, err = dataReader.ReadTable(); err != nil {
				return
			}
			if method.Mechanism, err = dataReader.ReadShortstr(); err != nil {
				return
			}

			if method.Response, err = dataReader.ReadLongstr(); err != nil {
				return
			}

			if method.Locale, err = dataReader.ReadShortstr(); err != nil {
				return
			}

		//case 20: // connection secure
		//case 21: // connection secure-ok
		//case 30: // connection tune
		case 31: // connection tune-ok
			method := connectionTuneOk{}
			mf.Method = &method
			if err = binary.Read(reader, binary.BigEndian, &method.ChannelMax); err != nil {
				return
			}

			if err = binary.Read(reader, binary.BigEndian, &method.FrameMax); err != nil {
				return
			}

			if err = binary.Read(reader, binary.BigEndian, &method.Heartbeat); err != nil {
				return
			}

		case 40: // connection open
			method := connectionOpen{}
			mf.Method = &method
			if method.VirtualHost, err = dataReader.ReadShortstr(); err != nil {
				return
			}

			if method.reserved1, err = dataReader.ReadShortstr(); err != nil {
				return
			}

			var bits byte
			if err = binary.Read(reader, binary.BigEndian, &bits); err != nil {
				return
			}
			method.reserved2 = bits&(1<<0) > 0

		//case 41: // connection open-ok
		//case 50: // connection close
		//case 51: // connection close-ok
		//case 60: // connection blocked
		//case 61: // connection unblocked
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 20: // channel
		switch methodId {
		//case 10: // channel open
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 40: // exchange
		switch methodId {
		//case 10: // exchange declare
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 50: // queue
		switch methodId {
		//case 10: // queue declare
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 60: // basic
		switch methodId {
		//case 10: // basic qos
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 85: // confirm
		switch methodId {
		//case 10: // confirm select
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	case 90: // tx
		switch methodId {
		//case 10: // tx select
		default:
			return mf, fmt.Errorf("bad method frame, unknown method %d for class %d", methodId, classId)
		}
	default:
		return mf, fmt.Errorf("bad method frame, unknown class %d", classId)
	}

	mf.ClassId = classId
	mf.MethodId = methodId

	return
}
