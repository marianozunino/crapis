package internal

import (
	"net"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
)

type Server struct {
	port       string
	bind       string
	listenAddr string
}

type ServerOption func(*Server)

func WithPort(port string) ServerOption {
	return func(s *Server) {
		log.Debug().Msgf("Configuring port to %s", port)
		s.port = port
	}
}

func WithBind(bind string) ServerOption {
	return func(s *Server) {
		log.Debug().Msgf("Configuring bind to %s", bind)
		s.bind = bind
	}
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		port: "6379",
		bind: "0.0.0.0",
	}

	for _, opt := range opts {
		opt(s)
	}

	s.listenAddr = net.JoinHostPort(s.bind, s.port)

	return s
}

func (s *Server) Run() {
	log.Debug().Msgf("Starting server on %s", s.listenAddr)
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer l.Close()

	// Update the listenAddr with the actual port if it was 0
	s.listenAddr = l.Addr().String()
	log.Debug().Msgf("Server is listening on %s", s.listenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %v", err)
			continue
		}
		go resp.HandleConnection(conn)
	}
}
