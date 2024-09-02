package resp

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

func (r *Reader) readInteger() (x int, n int, err error) {
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
