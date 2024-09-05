package server

import (
	"errors"
	"io"
	"net"
	"strings"

	"github.com/marianozunino/crapis/internal/command"
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
)

func (s *Server) handleConnection(conn net.Conn) {
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
			respWriter.Write(resp.NewError(err.Error()))
			continue
		}

		s.WriteToAOF(cmd, value)

		result := s.config.CmdExecutor.Execute(cmd, args)
		respWriter.Write(result)
	}
}

func (s *Server) WriteToAOF(cmd command.CommandType, value resp.Value) {
	if s.config.Aof == nil {
		return
	}

	if cmd == command.SET || cmd == command.SETEX || cmd == command.DEL {
		s.config.Aof.Write(value)
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

func (s *Server) loadAOF() error {
	if s.config.Aof == nil {
		return nil
	}

	log.Debug().Msg("Loading AOF")

	s.config.Aof.Read(func(r io.Reader) {
		respReader := resp.NewReader(r)
		for {
			value, err := respReader.Read()
			if err != nil {
				log.Debug().Msgf("Error reading from client: %s", err.Error())
				return
			}
			cmd, args, err := parseRequest(value)
			if err != nil {
				log.Debug().Msgf("Error reading from client: %s", err.Error())
				return
			}

			s.config.CmdExecutor.Execute(cmd, args)
		}
	})

	log.Debug().Msg("AOF loaded")

	return nil
}
