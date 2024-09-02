package resp

import (
	"bufio"
	"fmt"
	"io"
)

type Value struct {
	kind string

	strVal   string
	numVal   int
	bulkVal  string
	arrayVal []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) Resp {
	return Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) Read() (Value, error) {
	kind, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch kind {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(kind))
		return Value{}, fmt.Errorf("unknown type: %v", string(kind))
	}
}
