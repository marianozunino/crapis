package command

import (
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

type CommandHandler func([]resp.Value) resp.Value

type Executor struct {
	handlers map[CommandType]CommandHandler
	db       *store.Store // Assuming you have a Store type
}

func NewExecutor(db *store.Store) Executor {
	e := Executor{
		handlers: make(map[CommandType]CommandHandler),
		db:       db,
	}
	e.registerHandlers()
	return e
}

func (e *Executor) registerHandlers() {
	e.handlers[PING] = e.ping
	e.handlers[GET] = e.get
	e.handlers[SET] = e.set
	e.handlers[SETEX] = e.setex
	e.handlers[DEL] = e.del
}

func (e *Executor) Execute(cmd CommandType, args []resp.Value) resp.Value {
	handler, ok := e.handlers[cmd]
	if !ok {
		return resp.Value{
			Kind:   resp.ERROR,
			StrVal: "Unknown command",
		}
	}
	return handler(args)
}
