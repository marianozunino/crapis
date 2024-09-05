package resp

import (
	"reflect"
	"testing"
)

func TestReader_Marshal(t *testing.T) {
	tests := []struct {
		name  string
		input Value
		want  []byte
	}{
		{
			name:  "Marshal String",
			input: NewBulk(stringPtr("hello")),
			want:  []byte("+hello\r\n"),
		},

		{
			name:  "Marshal Bulk",
			input: NewBulk(stringPtr("hello")),
			want:  []byte("$5\r\nhello\r\n"),
		},

		{
			name:  "Marshal Array",
			input: NewArray(NewBulk(stringPtr("hello")), NewBulk(stringPtr("world"))),
			want:  []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
		},

		{
			name:  "Marshal Error",
			input: NewError("Error Message"),
			want:  []byte("-Error Message\r\n"),
		},

		{
			name:  "Marshal Positive Integer",
			input: NewInteger(42),
			want:  []byte(":42\r\n"),
		},
		{
			name:  "Marshal Negative Integer",
			input: NewInteger(-42),
			want:  []byte(":-42\r\n"),
		},
		{
			name:  "Marshal Null (RESP3)",
			input: NewNull(),
			want:  []byte("_\r\n"),
		},
		{
			name:  "Default",
			input: Value{},
			want:  []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Marshal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value.marshalString() = %s, want %s", got, tt.want)
			}
		})
	}

}
