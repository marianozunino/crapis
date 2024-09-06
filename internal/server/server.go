package server

import (
	"errors"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type Server struct {
	config *Config
	listen net.Listener
}

func NewServer(opts ...Option) *Server {
	config := newConfig(opts...)
	return &Server{config: config}
}

func (s *Server) Run() error {
	log.Info().Msgf("Starting server on %s", s.config.ListenAddr)
	l, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		return err
	}
	s.listen = l
	defer l.Close()

	s.config.ListenAddr = l.Addr().String()

	if err := s.loadAOF(); err != nil {
		log.Error().Err(err).Msg("Failed to load AOF")
		return err
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil // normal shutdown
			}
			log.Error().Err(err).Msg("Error accepting connection")
			continue
		}
		s.handleConnection(conn)
	}
}
