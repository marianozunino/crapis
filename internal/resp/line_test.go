package resp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestResp_readLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLine []byte
		wantN    int
		wantErr  bool
	}{
		{
			name:     "Valid line",
			input:    "hello\r\n",
			wantLine: []byte("hello"),
			wantN:    7,
			wantErr:  false,
		},
		{
			name:     "Empty line",
			input:    "\r\n",
			wantLine: []byte{},
			wantN:    2,
			wantErr:  false,
		},
		{
			name:     "Line without CRLF",
			input:    "hello",
			wantLine: nil,
			wantN:    0,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resp{
				reader: bufio.NewReader(strings.NewReader(tt.input)),
			}
			gotLine, gotN, err := r.readLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resp.readLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLine, tt.wantLine) {
				t.Errorf("Resp.readLine() gotLine = %v, want %v", gotLine, tt.wantLine)
			}
			if gotN != tt.wantN {
				t.Errorf("Resp.readLine() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
