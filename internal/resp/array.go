package resp

import "github.com/rs/zerolog/log"

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
