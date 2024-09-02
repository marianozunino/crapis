package resp

import "github.com/rs/zerolog/log"

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
