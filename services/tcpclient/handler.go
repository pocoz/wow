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

	hc := &hashcash.Hashcash{}
	err = decoder.Decode(&hc)
	if err != nil {
		level.Error(s.logger).Log("msg", "get first hc failure", "err", err)
		return err
	}

	hc.GeneratePow()

	err = encoder.Encode(hc)
	if err != nil {
		level.Error(s.logger).Log("msg", "encode pow hc failure", "err", err)
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
