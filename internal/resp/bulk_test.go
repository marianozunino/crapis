package resp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestReader_readBulk(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:    "Valid bulk string",
			input:   "$5\r\nhello\r\n",
			want:    Value{kind: BULK, bulkVal: stringPtr("hello")},
			wantErr: false,
		},
		{
			name:    "Empty bulk string",
			input:   "$0\r\n\r\n",
			want:    Value{kind: BULK, bulkVal: stringPtr("")},
			wantErr: false,
		},
		{
			name:    "Invalid bulk string size",
			input:   "$invalid\r\n",
			want:    Value{kind: BULK, bulkVal: nil},
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
				t.Errorf("Reader.readBulk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.readBulk() = %v, want %v", got, tt.want)
			}
		})
	}
}
