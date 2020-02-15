package data_transfer

type Table map[string]interface{}

// Decimal matches the AMQP decimal type.  Scale is the number of decimal
// digits Scale == 2, Value == 12345, Decimal == 123.45
type Decimal struct {
	Scale uint8
	Value int32
}
