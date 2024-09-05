package server

import (
	"net"

	"github.com/marianozunino/crapis/internal/aof"
	"github.com/marianozunino/crapis/internal/command"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port        string
	Bind        string
	ListenAddr  string
	CmdExecutor command.Executor
	Aof         *aof.Aof
}

type Option func(*Config)

func WithPort(port string) Option {
	return func(c *Config) {
		log.Debug().Msgf("Configuring port to %s", port)
		c.Port = port
	}
}

func WithBind(bind string) Option {
	return func(c *Config) {
		log.Debug().Msgf("Configuring bind to %s", bind)
		c.Bind = bind
	}
}

func WithCommandExecutor(executor command.Executor) Option {
	return func(c *Config) {
		c.CmdExecutor = executor
	}
}

func WithAof(aof *aof.Aof) Option {
	return func(c *Config) {
		c.Aof = aof
	}
}

func NewConfig(opts ...Option) *Config {
	c := &Config{
		Port: "6379",
		Bind: "0.0.0.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	c.ListenAddr = net.JoinHostPort(c.Bind, c.Port)
	return c
}
