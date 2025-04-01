package app

import (
	"context"
	"go.uber.org/zap"
	"order-service/internal/app/grpcapp"
	"order-service/internal/app/httpapp"
	"order-service/internal/config"
	"order-service/internal/logger"
	"order-service/internal/repository"
	"order-service/internal/storage"
	"order-service/internal/storage/postgresql"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpapp.App
	Storage    storage.Interface
}

func New(cfg *config.Config) *App {
	//strg := storage.New()
	s, err := postgresql.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	repo := repository.New(s)
	grpcApp := grpcapp.New(cfg.GPRCport, repo)
	httpApp := httpapp.New(cfg.HTTPport, cfg.GPRCport)

	return &App{
		GRPCServer: grpcApp,
		HTTPServer: httpApp,
		Storage:    s,
	}
}

func (a *App) Stop(ctx context.Context) {
	if err := a.HTTPServer.Stop(ctx); err != nil {
		logger.GetLogger().Error(ctx, "failed to stop HTTP server", zap.Error(err))
	}
	a.GRPCServer.Stop()
	a.Storage.Stop()
}
