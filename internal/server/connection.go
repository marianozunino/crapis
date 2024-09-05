package server

import (
	"errors"
	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, executor command.Executor) {
	defer conn.Close()
	respReader := resp.NewReader(conn)
	respWriter := resp.NewWriter(conn)
	for {
		value, err := respReader.Read()
		if err != nil {
			log.Debug().Msgf("Error reading from client: %s", err.Error())
			return
		}
		cmd, args, err := parseRequest(value)
		if err != nil {
			respWriter.Write(resp.Value{Kind: resp.ERROR, StrVal: err.Error()})
			continue
		}
		result := executor.Execute(cmd, args)
		respWriter.Write(result)
	}
}

func parseRequest(value resp.Value) (command.CommandType, []resp.Value, error) {
	if value.Kind != resp.ARRAY {
		return "", nil, errors.New("invalid request: expected array")
	}

	if len(value.ArrayVal) == 0 {
		return "", nil, errors.New("invalid request: empty array")
	}

	cmdValue := value.ArrayVal[0]
	if cmdValue.Kind != resp.BULK && cmdValue.Kind != resp.STRING {
		return "", nil, errors.New("invalid request: first array element must be bulk string or simple string")
	}

	var cmdStr string
	if cmdValue.Kind == resp.BULK {
		if cmdValue.BulkVal == nil {
			return "", nil, errors.New("invalid request: nil command")
		}
		cmdStr = *cmdValue.BulkVal
	} else {
		cmdStr = cmdValue.StrVal
	}

	cmd := command.CommandType(strings.ToUpper(cmdStr))
	args := value.ArrayVal[1:]

	return cmd, args, nil
}
