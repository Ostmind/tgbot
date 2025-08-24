package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	application "github.com/Ostmind/tgbot/internal/app"
	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/logger"
)

func main() {
	cfg := config.MustNew()

	sloger := logger.SetupLogger(cfg.Srv.EnvType)

	sloger.Info("starting Telegram Bot ToDoList")

	app, err := application.New(sloger, cfg)
	if err != nil {
		log.Fatal("No App cannot start server", slog.String("err", err.Error()))
	}

	pprofSrv := &http.Server{
		Addr:    ":6060",
		Handler: http.DefaultServeMux,
	}

	errCh := make(chan error, 1)

	go func() {
		pprofPort := "6060"
		sloger.Info("pprof server started", "port", pprofPort)
		if err := pprofSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	app.Run(cfg.Srv.Port)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	sloger.Info("Received interrupt signal")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Srv.ServerShutdownTimeout)
	defer cancel()

	app.Stop(ctx)
}
