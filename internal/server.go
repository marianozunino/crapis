package internal

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	port       string
	bind       string
	listenAddr string
}

type ServerOption func(*Server)

func WithPort(port string) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithBind(bind string) ServerOption {
	return func(s *Server) {
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
	fmt.Println("Listening on " + s.listenAddr)

	// Create a new server
	l, err := net.Listen("tcp", s.listenAddr)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()

	if err != nil {
		fmt.Println(err)
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
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}
