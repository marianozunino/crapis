package resp

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewResp(t *testing.T) {
	tests := []struct {
		name string
		args io.Reader
		want Resp
	}{
		{
			name: "Create new Resp",
			args: strings.NewReader("test"),
			want: Resp{reader: bufio.NewReader(strings.NewReader("test"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewResp(tt.args)
			if !reflect.DeepEqual(got.reader, tt.want.reader) {
				t.Errorf("NewResp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResp_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:    "Bulk string",
			input:   "$5\r\nhello\r\n",
			want:    Value{kind: "bulk", bulkVal: "hello"},
			wantErr: false,
		},
		{
			name:  "Array",
			input: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want: Value{kind: "array", arrayVal: []Value{
				{kind: "bulk", bulkVal: "hello"},
				{kind: "bulk", bulkVal: "world"},
			}},
			wantErr: false,
		},
		{
			name:    "Invalid kind",
			input:   "invalid\r\n",
			want:    Value{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resp{
				reader: bufio.NewReader(strings.NewReader(tt.input)),
			}
			got, err := r.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resp.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resp.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
