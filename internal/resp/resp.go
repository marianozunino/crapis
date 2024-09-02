package resp

import (
	"net"

	"github.com/rs/zerolog/log"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	for {

		respReader := NewReader(conn)

		value, err := respReader.Read()
		if err != nil {
			log.Debug().Msgf("Error reading from client: %s", err.Error())
			return
		}

		log.Debug().Msgf("Value: %+v", value)

		respWriter := NewWriter(conn)
		respWriter.Write(Value{kind: STRING, strVal: "OK"})
	}
}
