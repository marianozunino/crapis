package resp

// Define Redis type constants as a custom type.
type RedisType byte

const (
	// Redis protocol type indicators as per the Redis protocol specification.
	STRING  RedisType = '+'
	ERROR   RedisType = '-'
	INTEGER RedisType = ':'
	BULK    RedisType = '$'
	ARRAY   RedisType = '*'
	NULL    RedisType = '_'
)

const (
	// Carriage return and newline characters.
	CR byte = '\r'
	LF byte = '\n'
)
