package resp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestResp_readArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:  "Valid array",
			input: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want: Value{kind: "array", arrayVal: []Value{
				{kind: "bulk", bulkVal: "hello"},
				{kind: "bulk", bulkVal: "world"},
			}},
			wantErr: false,
		},
		{
			name:    "Empty array",
			input:   "*0\r\n",
			want:    Value{kind: "array", arrayVal: []Value{}},
			wantErr: false,
		},
		{
			name:    "Invalid array size",
			input:   "*invalid\r\n",
			want:    Value{kind: "array"},
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
				t.Errorf("Resp.readArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resp.readArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
