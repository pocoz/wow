package tcpserver

import (
	"fmt"
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/wow/db/local"
	"github.com/pocoz/wow/models"
)

type Service struct {
	logger   log.Logger
	listener net.Listener
	storage  *local.Storage
}

type Config struct {
	Logger  log.Logger
	Connect models.ConnectCfg
	Storage *local.Storage
}

func New(cfg *Config) (*Service, error) {
	listener, err := net.Listen(cfg.Connect.NetworkType, fmt.Sprintf("%s:%d", cfg.Connect.Address, cfg.Connect.Port))
	if err != nil {
		return nil, err
	}

	svc := &Service{
		logger:   cfg.Logger,
		listener: listener,
		storage:  cfg.Storage,
	}

	return svc, nil
}

func (s *Service) Run() error {
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			level.Info(s.logger).Log("msg", "tcpserver: error accepting connection", "err:", err)
			continue
		}

		level.Info(s.logger).Log("msg", "tcp accepted connection", "addr", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Service) Shutdown() {
	s.listener = nil

	level.Info(s.logger).Log("msg", "tcpserver: shutdown complete")
}
