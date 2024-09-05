package resp

import (
	"bytes"
	"fmt"
	"io"
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
		name    string
		args    Value
		want    string
		wantErr bool
	}{
		{
			name:    "Write Value",
			args:    Value{Kind: STRING, StrVal: "OK"},
			want:    "+OK\r\n",
			wantErr: false,
		},

		{
			name:    "Fail to write to the writer",
			args:    Value{Kind: STRING, StrVal: "OK"}, // Adjusted args to match the failing case
			want:    "",                                // Expected output is an empty string because of the error
			wantErr: true,                              // Expect an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b *bytes.Buffer
			var w io.Writer

			if tt.wantErr {
				// Use custom errorWriter to simulate a write error
				w = &errorWriter{}
			} else {
				// Use bytes.Buffer for successful writes
				b = &bytes.Buffer{}
				w = b
			}

			wr := &Writer{writer: w}
			err := wr.Write(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Only check the buffer's content if there's no error
				if got := b.String(); got != tt.want {
					t.Errorf("Writer.Write() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

type errorWriter struct{}

func (ew *errorWriter) Write(p []byte) (int, error) {
	// Return a custom error to simulate a failure in Write
	return 0, fmt.Errorf("write error")
}
