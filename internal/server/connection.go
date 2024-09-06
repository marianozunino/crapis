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
		if err := s.processCommand(respReader, respWriter); err != nil {
			if err != io.EOF {
				log.Debug().Err(err).Msg("Error processing command")
			}
			return
		}
	}
}

func (s *Server) processCommand(respReader resp.Reader, respWriter resp.Writer) error {
	value, err := respReader.Read()
	if err != nil {
		return err
	}

	cmd, args, err := parseRequest(value)
	if err != nil {
		return respWriter.Write(resp.NewError(err.Error()))
	}

	if err := s.WriteToAOF(cmd, value); err != nil {
		log.Error().Err(err).Msg("Failed to write to AOF")
	}

	result := s.config.CmdExecutor.Execute(cmd, args)
	return respWriter.Write(result)
}

func (s *Server) WriteToAOF(cmd command.CommandType, value resp.Value) error {
	if s.config.Aof == nil {
		return nil
	}
	if s.shouldWriteToAOF(cmd) {
		return s.config.Aof.Write(value)
	}
	return nil
}

func (s *Server) shouldWriteToAOF(cmd command.CommandType) bool {
	switch cmd {
	case command.SET:
		return true
	default:
		return false
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
	cmdStr := getCmdString(cmdValue)
	cmd := command.CommandType(strings.ToUpper(cmdStr))
	args := value.ArrayVal[1:]
	return cmd, args, nil
}

func getCmdString(cmdValue resp.Value) string {
	if cmdValue.Kind == resp.BULK {
		if cmdValue.BulkVal == nil {
			return ""
		}
		return *cmdValue.BulkVal
	}
	return cmdValue.StrVal
}

func (s *Server) loadAOF() error {
	if s.config.Aof == nil {
		return nil
	}
	log.Debug().Msg("Loading AOF")
	err := s.config.Aof.Read(func(r io.Reader) error {
		respReader := resp.NewReader(r)
		for {
			value, err := respReader.Read()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			cmd, args, err := parseRequest(value)
			if err != nil {
				return err
			}
			s.config.CmdExecutor.Execute(cmd, args)
		}
	})
	if err != nil {
		return err
	}
	log.Debug().Msg("AOF loaded")
	return nil
}

