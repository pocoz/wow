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

		blockForClient := hashcash.NewBlock(msg.Body)
		err = encoder.Encode(blockForClient)
		if err != nil {
			level.Error(s.logger).Log("msg", "send block failure", "err", err)
			return
		}

		blockFromClient := &hashcash.Block{}
		err = decoder.Decode(&blockFromClient)
		if err != nil {
			level.Error(s.logger).Log("msg", "get pow block failure", "err", err)
			return
		}

		if s.hcSvc.ValidateBlock(blockForClient, blockFromClient) {
			msg = &models.Message{
				Body: tools.GetQuote(),
			}

			err = encoder.Encode(msg)
			if err != nil {
				level.Error(s.logger).Log("msg", "send final message failure", "err", err)
			}
		} else {
			level.Info(s.logger).Log("msg", "block is not valid")
		}
	}
}
