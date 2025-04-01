package httpapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"order-service/pkg/api/test"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order-service/internal/logger"
)

type App struct {
	HTTPServer *http.Server
	HTTPport   int
	GRPCport   int
}

func New(HTTPport, GRPCport int) *App {
	return &App{
		HTTPport: HTTPport,
		GRPCport: GRPCport,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		logger.GetLogger().Fatal(context.Background(), err.Error())
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	mux := runtime.NewServeMux()
	a.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.HTTPport),
		Handler: mux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("0.0.0.0:%d", a.GRPCport)
	err := test.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, addr, opts)
	if err != nil {
		logger.GetLogger().Error(ctx, "failed to register gRPC gateway", zap.Error(err))
		return fmt.Errorf("%s: failed to register gateway: %w", op, err)
	}

	logger.GetLogger().Info(ctx, "http server started", zap.String("addr", a.HTTPServer.Addr))
	if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: failed to serve HTTP: %w", op, err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	logger.GetLogger().Info(context.Background(), "http server stopped")
	return nil
}
