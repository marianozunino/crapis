package resp

import (
	"reflect"
	"testing"
)

func TestReader_marshalString(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal string",
			input: Value{Kind: STRING, StrVal: "hello"},
			want:  []byte("+hello\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.marshalString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestReader_marshalBulk(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal Bulk",
			input: Value{Kind: BULK, BulkVal: stringPtr("hello")},
			want:  []byte("$5\r\nhello\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.marshalBulk(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestReader_marshalArray(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal Array",
			input: Value{Kind: ARRAY, ArrayVal: []Value{{Kind: BULK, BulkVal: stringPtr("hello")}, {Kind: BULK, BulkVal: stringPtr("world")}}},
			want:  []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.marshalArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestReader_marshallNull(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal Null",
			input: Value{Kind: NULL},
			want:  []byte("_\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.marshallNull(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestReader_marshallError(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal Error",
			input: Value{Kind: ERROR, StrVal: "Error Message"},
			want:  []byte("-Error Message\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.marshallError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}
}
