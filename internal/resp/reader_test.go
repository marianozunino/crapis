package resp

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	tests := []struct {
		name string
		args io.Reader
		want Reader
	}{
		{
			name: "Create new Reader",
			args: strings.NewReader("test"),
			want: Reader{reader: bufio.NewReader(strings.NewReader("test"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewReader(tt.args)
			if !reflect.DeepEqual(got.reader, tt.want.reader) {
				t.Errorf("NewReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:    "Bulk string",
			input:   "$5\r\nhello\r\n",
			want:    Value{Kind: BULK, BulkVal: stringPtr("hello")},
			wantErr: false,
		},
		{
			name:    "Array",
			input:   "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want:    NewArray(NewBulk(stringPtr("hello")), NewBulk(stringPtr("world"))),
			wantErr: false,
		},
		{
			name:    "Invalid kind",
			input:   "invalid\r\n",
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "Fail to read byte from reader",
			input:   "",
			want:    Value{},
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
				t.Errorf("Reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
