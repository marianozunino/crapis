package resp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) Reader {
	return Reader{
		reader: bufio.NewReader(rd),
	}
}

func (r *Reader) Read() (Value, error) {
	kind, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch RedisType(kind) {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		log.Debug().Msgf("Unknown type: %v", string(kind))
		return Value{}, fmt.Errorf("unknown type: %v", string(kind))
	}
}
