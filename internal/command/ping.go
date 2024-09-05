package command

import "github.com/marianozunino/crapis/internal/resp"

func (e *executor) ping(args []resp.Value) resp.Value {
	if len(args) > 1 {
		return resp.NewError("wrong number of arguments for 'ping' command")
	}

	if len(args) == 1 {
		return resp.NewString(*args[0].BulkVal)
	}

	return resp.NewString("PONG")
}
