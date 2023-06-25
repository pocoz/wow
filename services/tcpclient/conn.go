package tcpclient

import (
	"fmt"
	"net"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/wow/models"
)

type Service struct {
	logger log.Logger
	conn   net.Conn
}

type Config struct {
	Logger  log.Logger
	Connect models.ConnectCfg
}

func New(cfg *Config) (*Service, error) {
	conn, err := net.Dial(cfg.Connect.NetworkType, fmt.Sprintf("%s:%d", cfg.Connect.Address, cfg.Connect.Port))
	if err != nil {
		return nil, err
	}

	svc := &Service{
		logger: cfg.Logger,
		conn:   conn,
	}

	return svc, nil
}

func (s *Service) Run() error {
	defer s.conn.Close()
	for {
		err := s.handleConnection()
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Second)
	}
}

func (s *Service) Shutdown() {
	s.conn.Close()

	level.Info(s.logger).Log("msg", "tcpclient: shutdown complete")
}
