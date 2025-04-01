package grpcapp

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"order-service/internal/logger"
	"order-service/internal/middleware"
	"order-service/internal/service"
)

type App struct {
	GRPCServer *grpc.Server
	port       int
}

func New(port int, repo service.OrderRepository) *App {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.LoggerInterceptor))
	service.Register(grpcServer, repo)

	return &App{
		GRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		logger.GetLogger().Fatal(context.Background(), err.Error())
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	logger.GetLogger().Info(context.Background(), "grpc server started", zap.String("addr", lis.Addr().String()))
	if err := a.GRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {

	logger.GetLogger().Info(context.Background(), "grpc server stopped")
	a.GRPCServer.GracefulStop()
}
