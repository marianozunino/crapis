package server

import (
	"net"

	"github.com/rs/zerolog/log"
)

type Server struct {
	config *Config
}

func NewServer(config *Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run() error {
	log.Info().Msgf("Starting server on %s", s.config.ListenAddr)
	l, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	defer l.Close()
	s.config.ListenAddr = l.Addr().String()

	s.loadAOF()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Error().Msgf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}
