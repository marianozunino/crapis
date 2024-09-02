package resp

import (
	"net"
	"strings"

	"github.com/marianozunino/crapis/internal/resp/command"
	"github.com/rs/zerolog/log"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	for {

		respReader := NewReader(conn)
		respWriter := NewWriter(conn)

		value, err := respReader.Read()
		if err != nil {
			log.Debug().Msgf("Error reading from client: %s", err.Error())
			return
		}

		if value.kind != ARRAY {
			log.Error().Msg("Invalid request, expected array")
			continue
		}

		if len(value.arrayVal) == 0 {
			log.Error().Msg("Invalid request, expected array length > 0")
			continue
		}

		cmdStr := strings.ToUpper(*value.arrayVal[0].bulkVal)

		cmd, err := command.ParseCommand(cmdStr)

		if err != nil {
			respWriter.Write(Value{kind: ERROR, strVal: err.Error()})
			continue
		}

		args := value.arrayVal[1:]

		log.Debug().Msgf("Value: %+v", value)

		handler, ok := Handlers[cmd]
		if !ok {
			log.Error().Msgf("Invalid command %s", err)
			respWriter.Write(Value{kind: ERROR, strVal: err.Error()})
			continue
		}

		result := handler(args)

		respWriter.Write(result)
	}
}
