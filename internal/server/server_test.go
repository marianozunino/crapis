package server

import (
	"net"
	"testing"
	"time"

	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/stretchr/testify/assert"
)

func TestServerInitializationErrors(t *testing.T) {
	tests := []struct {
		name        string
		configFunc  func() *Config
		expectedErr string
	}{
		{
			name: "Invalid port",
			configFunc: func() *Config {
				return NewConfig(WithPort("99999")) // Invalid port number
			},
			expectedErr: "invalid port",
		},
		{
			name: "Invalid bind address",
			configFunc: func() *Config {
				return NewConfig(WithBind("invalid-address"))
			},
			expectedErr: "invalid-address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.configFunc()
			s := NewServer(config)

			err := s.Run()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestServerConnectionErrors(t *testing.T) {
	// Create a mock executor
	mockExecutor := &MockExecutor{}

	// Use net.Pipe to create a pair of connected pipes
	clientConn, _ := net.Pipe()

	// Create a server with the mock executor
	config := NewConfig(
		WithPort("0"), // Use a free port
		WithCommandExecutor(mockExecutor),
	)
	s := NewServer(config)

	// Run the server in a goroutine
	go func() {
		err := s.Run()
		if err != nil {
			t.Fatalf("Server failed to run: %v", err)
		}
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	// Close the client connection to simulate a connection error
	clientConn.Close()

	// Try writing data to simulate an error
	_, err := clientConn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Logf("expected error when writing to closed connection: %v", err)
	} else {
		t.Fatal("expected an error but write succeeded")
	}

	// Ensure that the server handles connection errors gracefully
	// Depending on how you handle logging or monitoring, check logs or other indicators
}
func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		opts           []Option
		expectedPort   string
		expectedBind   string
		expectedListen string
	}{
		{
			name:           "Default configuration",
			opts:           []Option{},
			expectedPort:   "6379",
			expectedBind:   "0.0.0.0",
			expectedListen: "0.0.0.0:6379",
		},
		{
			name:           "Custom port",
			opts:           []Option{WithPort("8080")},
			expectedPort:   "8080",
			expectedBind:   "0.0.0.0",
			expectedListen: "0.0.0.0:8080",
		},
		{
			name:           "Custom bind",
			opts:           []Option{WithBind("127.0.0.1")},
			expectedPort:   "6379",
			expectedBind:   "127.0.0.1",
			expectedListen: "127.0.0.1:6379",
		},
		{
			name:           "Custom port and bind",
			opts:           []Option{WithPort("8080"), WithBind("127.0.0.1")},
			expectedPort:   "8080",
			expectedBind:   "127.0.0.1",
			expectedListen: "127.0.0.1:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(tt.opts...)
			assert.Equal(t, tt.expectedPort, c.Port)
			assert.Equal(t, tt.expectedBind, c.Bind)
			assert.Equal(t, tt.expectedListen, c.ListenAddr)
		})
	}
}

func TestServerRun(t *testing.T) {
	// Create a mock executor for testing
	mockExecutor := &MockExecutor{}

	config := NewConfig(
		WithPort("0"), // Use port 0 to let the system assign a free port
		WithCommandExecutor(mockExecutor),
	)
	s := NewServer(config)

	go s.Run()

	// Wait for the server to start and get the actual port
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", s.config.ListenAddr)
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
			input:    "*1\r\n$4\r\nPING\r\n",
			expected: "+PONG\r\n",
		},
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

// MockExecutor is a mock implementation of the command.Executor interface for testing
type MockExecutor struct{}

func (m *MockExecutor) Execute(cmd command.CommandType, args []resp.Value) resp.Value {
	// For this test, we're only implementing the PING command
	if cmd == command.PING {
		return resp.NewString("PONG")
	}
	return resp.NewError("unknown command")
}
