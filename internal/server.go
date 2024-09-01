package internal

import (
	"io"
	"net"
	"os"

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

	// Create a new server
	l, err := net.Listen("tcp", s.listenAddr)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// Listen for connections
	conn, err := l.Accept()

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Debug().Msgf("error reading from client: %s", err.Error())
			os.Exit(1)
		}
		log.Debug().Msgf("message from client: %s", string(buf))

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}
