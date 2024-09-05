package command

import (
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

type CommandHandler func([]resp.Value) resp.Value

type executor struct {
	handlers map[CommandType]CommandHandler
	db       *store.Store // Assuming you have a Store type
}

// define a Executor interface
type Executor interface {
	Execute(cmd CommandType, args []resp.Value) resp.Value
}

func NewExecutor(db *store.Store) Executor {
	e := executor{
		handlers: make(map[CommandType]CommandHandler),
		db:       db,
	}
	e.registerHandlers()
	return &e
}

func (e *executor) registerHandlers() {
	e.handlers[PING] = e.ping
	e.handlers[GET] = e.get
	e.handlers[SET] = e.set
	e.handlers[SETEX] = e.setex
	e.handlers[DEL] = e.del
	e.handlers[EXPIRE] = e.expire
}

func (e *executor) Execute(cmd CommandType, args []resp.Value) resp.Value {
	handler, ok := e.handlers[cmd]
	if !ok {
		return resp.NewError("Unknown command")
	}
	return handler(args)
}
