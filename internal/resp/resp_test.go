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

func TestResp_readBulk(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Value
		wantErr bool
	}{
		{
			name:    "Valid bulk string",
			input:   "$5\r\nhello\r\n",
			want:    Value{kind: "bulk", bulkVal: "hello"},
			wantErr: false,
		},
		{
			name:    "Empty bulk string",
			input:   "$0\r\n\r\n",
			want:    Value{kind: "bulk", bulkVal: ""},
			wantErr: false,
		},
		{
			name:    "Invalid bulk string size",
			input:   "$invalid\r\n",
			want:    Value{kind: "bulk"},
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
				t.Errorf("Resp.readBulk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resp.readBulk() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
