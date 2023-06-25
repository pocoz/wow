package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log/level"
	"github.com/pocoz/wow/models"
	"github.com/pocoz/wow/services/hashcash"
)

func (s *Service) handleConnection() error {
	decoder := json.NewDecoder(s.conn)
	encoder := json.NewEncoder(s.conn)

	msg := &models.Message{Body: s.conn.LocalAddr().String()}
	err := encoder.Encode(msg)
	if err != nil {
		level.Error(s.logger).Log("msg", "send first message failure", "err", err)
		return err
	}

	block := &hashcash.Block{}
	err = decoder.Decode(&block)
	if err != nil {
		level.Error(s.logger).Log("msg", "get first block failure", "err", err)
		return err
	}

	block.GeneratePow()

	err = encoder.Encode(block)
	if err != nil {
		level.Error(s.logger).Log("msg", "encode pow block failure", "err", err)
		return err
	}

	err = decoder.Decode(&msg)
	if err != nil {
		level.Error(s.logger).Log("msg", "get final message failure", "err", err)
		return err
	}

	fmt.Println(msg.Body)

	return nil
}
