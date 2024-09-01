package internal

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name           string
		opts           []ServerOption
		expectedPort   string
		expectedBind   string
		expectedListen string
	}{
		{
			name:           "Default configuration",
			opts:           []ServerOption{},
			expectedPort:   "6379",
			expectedBind:   "0.0.0.0",
			expectedListen: "0.0.0.0:6379",
		},
		{
			name:           "Custom port",
			opts:           []ServerOption{WithPort("8080")},
			expectedPort:   "8080",
			expectedBind:   "0.0.0.0",
			expectedListen: "0.0.0.0:8080",
		},
		{
			name:           "Custom bind",
			opts:           []ServerOption{WithBind("127.0.0.1")},
			expectedPort:   "6379",
			expectedBind:   "127.0.0.1",
			expectedListen: "127.0.0.1:6379",
		},
		{
			name:           "Custom port and bind",
			opts:           []ServerOption{WithPort("8080"), WithBind("127.0.0.1")},
			expectedPort:   "8080",
			expectedBind:   "127.0.0.1",
			expectedListen: "127.0.0.1:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.opts...)
			assert.Equal(t, tt.expectedPort, s.port)
			assert.Equal(t, tt.expectedBind, s.bind)
			assert.Equal(t, tt.expectedListen, s.listenAddr)
		})
	}
}

func TestServerRun(t *testing.T) {
	s := NewServer(WithPort("0")) // Use port 0 to let the system assign a free port

	go s.Run()

	// Wait for the server to start and get the actual port
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", s.listenAddr)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "PING command",
			input:    "PING\r\n",
			expected: "+OK\r\n",
		},
		// Add more test cases here if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conn.Write([]byte(tt.input))
			assert.NoError(t, err)

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(buf[:n]))
		})
	}
}

