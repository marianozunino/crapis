package resp

func (r *Reader) readBulk() (Value, error) {
	v := Value{}

	v.kind = BULK

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulkVal = new(string)
	*v.bulkVal = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
