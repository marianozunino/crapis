package resp

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
