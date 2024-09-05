package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewWriter(t *testing.T) {
	tests := []struct {
		name  string
		want  *Writer
		wantW string
	}{
		{
			name:  "Create new Writer",
			want:  &Writer{writer: &bytes.Buffer{}},
			wantW: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := NewWriter(w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewWriter() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {

	tests := []struct {
		name string
		args Value
		want string
	}{
		{
			name: "Write Value",
			args: Value{Kind: STRING, StrVal: "OK"},
			want: "+OK\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}

			w := &Writer{writer: b}
			w.Write(tt.args)

			if got := b.String(); got != tt.want {
				t.Errorf("Writer.Write() = %v, want %v", got, tt.want)
			}

		})
	}

}
