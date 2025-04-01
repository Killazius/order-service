package main

import (
	"context"
	"fmt"
	"order-service/internal/app"
	"order-service/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		os.Exit(1)
	}
	application := app.New(cfg)
	go application.GRPCServer.MustRun()
	go application.HTTPServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.Stop(ctx)

	// указано добавить по заданию
	fmt.Println("Server Stopped")
}
