package resp

import (
	"strconv"
)

// Value represents a parsed Redis value with different types.
type Value struct {
	Kind RedisType

	// Specific fields to store different types of Redis values.
	StrVal   string  // For simple strings and error messages.
	NumVal   int     // For integer values.
	BulkVal  *string // For bulk strings (nil if not present).
	ArrayVal []Value // For arrays.
}

func (v Value) Marshal() []byte {
	switch v.Kind {
	case ARRAY:
		return v.marshalArray()
	case BULK:
		return v.marshalBulk()
	case STRING:
		return v.marshalString()
	case INTEGER:
		return v.marshalInteger()
	case NULL:
		return v.marshallNull()
	case ERROR:
		return v.marshallError()
	default:
		return []byte{}
	}
}

// marshalString returns the Redis protocol representation of a simple string
// Example: "+hello\r\n"
func (v Value) marshalString() []byte {
	var bytes []byte

	bytes = append(bytes, byte(STRING))
	bytes = append(bytes, v.StrVal...)
	bytes = append(bytes, CR, LF)

	return bytes
}

// marshalBulk returns the Redis protocol representation of a bulk string
// Example: "$5\r\nhello\r\n"
func (v Value) marshalBulk() []byte {
	var bytes []byte

	bytes = append(bytes, byte(BULK))
	bytes = append(bytes, strconv.Itoa(len(*v.BulkVal))...)
	bytes = append(bytes, CR, LF)
	bytes = append(bytes, *v.BulkVal...)
	bytes = append(bytes, CR, LF)

	return bytes
}

// marshalArray returns the Redis protocol representation of an array
// Example: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
func (v Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, byte(ARRAY))
	bytes = append(bytes, strconv.Itoa(len(v.ArrayVal))...)
	bytes = append(bytes, CR, LF)

	for _, val := range v.ArrayVal {
		bytes = append(bytes, val.Marshal()...)
	}

	return bytes
}

// marshallNull returns the Redis protocol representation of a null value
// Example: "_\r\n"
func (v Value) marshallNull() []byte {
	var bytes []byte
	bytes = append(bytes, byte(NULL))
	bytes = append(bytes, CR, LF)
	return bytes
}

// marshallError returns the Redis protocol representation of an error message
// Example: "-Error message\r\n"
func (v Value) marshallError() []byte {

	var bytes []byte
	bytes = append(bytes, byte(ERROR))
	bytes = append(bytes, v.StrVal...)
	bytes = append(bytes, CR, LF)

	return bytes
}

// marshallInteger returns the Redis protocol representation of an integer
// Example: ":123\r\n"
func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, byte(INTEGER))

	if v.NumVal < 0 {
		bytes = append(bytes, '-')
	} else {
		bytes = append(bytes, '+')
	}

	bytes = append(bytes, strconv.Itoa(v.NumVal)...)

	bytes = append(bytes, CR, LF)
	return bytes
}
