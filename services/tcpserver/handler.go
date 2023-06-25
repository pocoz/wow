package tcpserver

import (
	"encoding/json"
	"net"

	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/wow/models"
	"github.com/pocoz/wow/services/hashcash"
	"github.com/pocoz/wow/tools"
)

func (s *Service) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg *models.Message
		err := decoder.Decode(&msg)
		if err != nil {
			level.Error(s.logger).Log("msg", "get first message failure", "err", err)
			return
		}

		hc := hashcash.NewHashcash(msg.Body)

		err = encoder.Encode(hc)
		if err != nil {
			level.Error(s.logger).Log("msg", "send hc failure", "err", err)
			return
		}

		err = decoder.Decode(&hc)
		if err != nil {
			level.Error(s.logger).Log("msg", "get pow hc failure", "err", err)
			return
		}

		if hc.Validate() {
			msg = &models.Message{
				Body: tools.GetQuote(),
			}

			err = encoder.Encode(msg)
			if err != nil {
				level.Error(s.logger).Log("msg", "send final message failure", "err", err)
			}
		} else {
			level.Info(s.logger).Log("msg", "hc is not valid")
		}
	}
}
