package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Ostmind/tgbot/internal/app"
	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/logger"
)

func main() {
	cfg := config.MustNew()

	sloger := logger.SetupLogger(cfg.Srv.EnvType)

	sloger.Info("starting Telegram Bot ToDoList")

	app, err := app.New(sloger, cfg)
	if err != nil {
		log.Fatal("No App cannot start server", slog.Any("error", err))
	}

	app.Run()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	sloger.Info("Received interrupt signal")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Srv.ServerShutdownTimeout)
	defer cancel()

	app.Stop(ctx)
}
