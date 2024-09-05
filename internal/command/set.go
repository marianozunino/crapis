package command

import (
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
)

func (e *Executor) set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewError("wrong number of arguments for 'set' command")
	}

	key := args[0].BulkVal
	val := args[1].BulkVal

	log.Debug().Msgf("set command [%s] = %s", *key, *val)

	e.db.StoreValue(*key, *val)

	return resp.Ok()
}
