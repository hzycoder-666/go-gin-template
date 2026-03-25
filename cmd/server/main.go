package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "hzycoder.com/go-gin-template/internal/config"
	"hzycoder.com/go-gin-template/internal/database"
	router "hzycoder.com/go-gin-template/internal/routes"
	logger "hzycoder.com/go-gin-template/pkg/logger"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase(config.Global.Database.DSN)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	logger.Init()

	r := router.SetupRouter()

	port := config.Global.Server.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		slog.Info("server start", "port", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
}
