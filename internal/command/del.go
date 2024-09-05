package command

import (
	"github.com/marianozunino/crapis/internal/resp"
)

func (e *Executor) del(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewError("wrong number of arguments for 'del' command")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = *arg.BulkVal
	}
	deletedKeys := e.db.DeleteKey(keys...)
	return resp.NewInteger(deletedKeys)
}
