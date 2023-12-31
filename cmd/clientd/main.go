package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/wow/config"
	"github.com/pocoz/wow/models"
	"github.com/pocoz/wow/services/tcpclient"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		exitCodeSuccess = 0
		exitCodeFailure = 1
	)

	donec := make(chan struct{})
	sigc := make(chan os.Signal, 1)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "app", config.ServiceName, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg, err := config.New()
	if err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(exitCodeFailure)
	}

	clientSvc, err := tcpclient.New(&tcpclient.Config{
		Logger: logger,
		Connect: models.ConnectCfg{
			NetworkType: cfg.ServerType,
			Address:     cfg.ServerAddress,
			Port:        cfg.ServerPort,
		},
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to tcp client", "err", err)
		os.Exit(exitCodeFailure)
	}
	go func() {
		level.Info(logger).Log("msg", "starting tcp client", "port", cfg.ServerPort)
		err = clientSvc.Run()
		if err != nil {
			level.Error(logger).Log("msg", "tcp client run failure", "err", err)
			os.Exit(exitCodeFailure)
		}
	}()

	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	go func() {
		select {
		case sig := <-sigc:
			level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)
			clientSvc.Shutdown()
			signal.Stop(sigc)
			close(donec)
		}
	}()

	<-donec
	level.Info(logger).Log("msg", "goodbye")
	os.Exit(exitCodeSuccess)
}
