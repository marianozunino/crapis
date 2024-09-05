package server

import (
	"github.com/rs/zerolog/log"
	"net"
)

type Server struct {
	config *Config
}

func NewServer(config *Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run() {
	log.Debug().Msgf("Starting server on %s", s.config.ListenAddr)
	l, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer l.Close()
	s.config.ListenAddr = l.Addr().String()
	log.Debug().Msgf("Server is listening on %s", s.config.ListenAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, s.config.CmdExecutor)
	}
}
