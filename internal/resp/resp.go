package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	// https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

const (
	R byte = '\r'
	N      = '\n'
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

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.kind = "array"

	// read length of array
	size, _, err := r.readInteger()
	if err != nil {
		log.Debug().Msgf("Error reading array length: %v", err)
		return v, err
	}

	// foreach line, parse and read the value
	v.arrayVal = make([]Value, 0)

	for i := 0; i < size; i++ {
		val, err := r.Read()

		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.arrayVal = append(v.arrayVal, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.kind = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulkVal = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == R {
			log.Debug().Msgf("Read line: %s", string(line))
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	log.Debug().Msgf("Read integer: %d", i64)

	return int(i64), n, nil
}
