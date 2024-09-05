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
			input: Value{Kind: STRING, StrVal: "hello"},
			want:  []byte("+hello\r\n"),
		},

		{
			name:  "Marshal Bulk",
			input: Value{Kind: BULK, BulkVal: stringPtr("hello")},
			want:  []byte("$5\r\nhello\r\n"),
		},

		{
			name: "Marshal Array",
			input: Value{Kind: ARRAY, ArrayVal: []Value{
				{Kind: BULK, BulkVal: stringPtr("hello")},
				{Kind: BULK, BulkVal: stringPtr("world")},
			}},
			want: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
		},

		{
			name:  "Marshal Error",
			input: Value{Kind: ERROR, StrVal: "Error Message"},
			want:  []byte("-Error Message\r\n"),
		},

		{
			name:  "Marshal Positive Integer",
			input: Value{Kind: INTEGER, NumVal: 42},
			want:  []byte(":+42\r\n"),
		},
		{
			name:  "Marshal Negative Integer",
			input: Value{Kind: INTEGER, NumVal: -42},
			want:  []byte(":-42\r\n"),
		},
		{
			name:  "Marshal Null (RESP3)",
			input: Value{Kind: NULL},
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
