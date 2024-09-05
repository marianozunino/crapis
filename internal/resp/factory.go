package resp

// Factory functions to create specific types of Redis Values

// NewError creates a new error Value
func NewError(message string) Value {
	return Value{
		Kind:   ERROR,
		StrVal: message,
	}
}

// NewString creates a new simple string Value
func NewString(message string) Value {
	return Value{
		Kind:   STRING,
		StrVal: message,
	}
}

// NewBulk creates a new bulk string Value
func NewBulk(message *string) Value {
	return Value{
		Kind:    BULK,
		BulkVal: message,
	}
}

// NewInteger creates a new integer Value
func NewInteger(number int) Value {
	return Value{
		Kind:   INTEGER,
		NumVal: number,
	}
}

// NewArray creates a new array Value
func NewArray(values ...Value) Value {

	return Value{
		Kind:     ARRAY,
		ArrayVal: values,
	}
}

// NewNull creates a new null Value
func NewNull() Value {
	return Value{
		Kind: NULL,
	}
}

// Ok creates a new simple string Value with the message "OK"
func Ok() Value {
	return Value{
		Kind:   STRING,
		StrVal: "OK",
	}
}
