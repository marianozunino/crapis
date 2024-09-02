package resp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func stringPtr(s string) *string {
	return &s
}

func TestReader_readArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:  "Valid array",
			input: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want: Value{kind: ARRAY, arrayVal: []Value{
				{kind: BULK, bulkVal: stringPtr("hello")},
				{kind: BULK, bulkVal: stringPtr("world")},
			}},
			wantErr: false,
		},
		{
			name:    "Empty array",
			input:   "*0\r\n",
			want:    Value{kind: ARRAY, arrayVal: []Value{}},
			wantErr: false,
		},
		{
			name:    "Invalid array size",
			input:   "*invalid\r\n",
			want:    Value{kind: ARRAY},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reader{
				reader: bufio.NewReader(strings.NewReader(tt.input)),
			}
			got, err := r.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.readArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.readArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
