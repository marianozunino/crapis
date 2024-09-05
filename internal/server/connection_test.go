package server

import (
	"testing"

	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/stretchr/testify/assert"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name         string
		value        resp.Value
		expectedCmd  command.CommandType
		expectedArgs []resp.Value
		expectedErr  bool
	}{
		{
			name:         "Simple string",
			value:        resp.NewArray(resp.NewString("PING")),
			expectedCmd:  command.PING,
			expectedArgs: []resp.Value{},
		},
		{
			name:         "Empty array",
			value:        resp.NewArray(),
			expectedArgs: []resp.Value{},
			expectedErr:  true,
		},
		{
			name:         "First arg is not a string or bulk string (aka command)",
			value:        resp.NewArray(resp.NewInteger(1)),
			expectedArgs: []resp.Value{},
			expectedErr:  true,
		},
		{
			name:         "Bulk arg has no value",
			value:        resp.NewArray(resp.NewBulk(nil)),
			expectedArgs: []resp.Value{},
			expectedErr:  true,
		},
		{
			name:         "Root node is not an array",
			value:        resp.NewString("PING"),
			expectedArgs: []resp.Value{},
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args, err := parseRequest(tt.value)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCmd, cmd)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}
