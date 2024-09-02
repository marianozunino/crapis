package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestResp_readInteger(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantX   int
		wantN   int
		wantErr bool
	}{
		{
			name:    "Positive integer",
			input:   "1000\r\n",
			wantX:   1000,
			wantN:   6,
			wantErr: false,
		},
		{
			name:    "Negative integer",
			input:   "-1000\r\n",
			wantX:   -1000,
			wantN:   7,
			wantErr: false,
		},
		{
			name:    "Zero",
			input:   "0\r\n",
			wantX:   0,
			wantN:   3,
			wantErr: false,
		},
		{
			name:    "Invalid integer",
			input:   "invalid\r\n",
			wantX:   0,
			wantN:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resp{
				reader: bufio.NewReader(strings.NewReader(tt.input)),
			}
			gotX, gotN, err := r.readInteger()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resp.readInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotX != tt.wantX {
				t.Errorf("Resp.readInteger() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotN != tt.wantN {
				t.Errorf("Resp.readInteger() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
