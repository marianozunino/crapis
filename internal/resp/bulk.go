package resp

func (r *Reader) readBulk() (Value, error) {
	v := Value{}

	v.Kind = BULK

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.BulkVal = new(string)
	*v.BulkVal = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
