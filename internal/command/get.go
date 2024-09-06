package command

import (
	"github.com/marianozunino/crapis/internal/resp"
	"github.com/rs/zerolog/log"
)

func (e *executor) get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("wrong number of arguments for 'get' command")
	}

	key := args[0].BulkVal

	log.Debug().Msgf("get command [%s]", *key)

	val := e.db.ReadValue(*key)

	if val == nil {
		return resp.NewNull()
	}

	return resp.NewString(*val)
}
