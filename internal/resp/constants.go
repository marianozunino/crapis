package resp

const (
	// https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings
	STRING  byte = '+'
	ERROR   byte = '-'
	INTEGER byte = ':'
	BULK    byte = '$'
	ARRAY   byte = '*'
)

const (
	R byte = '\r'
	N byte = '\n'
)
