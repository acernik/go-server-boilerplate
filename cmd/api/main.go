package main

import (
	"go.uber.org/zap"

	"github.com/acernik/go-server-boilerplate/internal/app"
	"github.com/acernik/go-server-boilerplate/internal/server"
	"github.com/acernik/go-server-boilerplate/internal/service"
	"github.com/acernik/go-server-boilerplate/internal/store"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := app.LoadConfig()
	if err != nil {
		return err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer logger.Sync() //nolint:errcheck

	db, err := store.NewStore(cfg)
	if err != nil {
		return err
	}

	svc := service.NewService(db)

	srv := server.New(svc, logger, cfg.App.TimeFormat)
	r := srv.RegisterRoutes()

	err = r.Run(cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
