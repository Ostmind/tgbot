package app

import (
	"context"
	"fmt"
	"github.com/Ostmind/tgbot/internal/bot"
	"log/slog"
	"os"

	"github.com/Ostmind/tgbot/internal/config"
	userhandler "github.com/Ostmind/tgbot/internal/server/handler/user"
	srv "github.com/Ostmind/tgbot/internal/server/server"
	"github.com/Ostmind/tgbot/internal/services/users"
	"github.com/Ostmind/tgbot/internal/storage/postgres"
)

type App struct {
	server *srv.Server
	logger *slog.Logger
	db     *postgres.Storage
	cfg    *config.AppConfig
}

func New(logger *slog.Logger, cfg *config.AppConfig) (*App, error) {
	db, err := postgres.New(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish db connection %w", err)
	}

	userManager := users.New(db)

	userHandler := userhandler.NewUserHandler(userManager, cfg, logger)

	server := srv.New(logger, cfg, db, userHandler)

	return &App{
		server: server,
		logger: logger,
		db:     db,
		cfg:    cfg,
	}, nil
}

func (a App) Run() {
	a.logger.Info("Starting app...")
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		panic("TELEGRAM_TOKEN not set")
	}

	bot.RunBot(token, a.db)

	a.server.Run()
}

func (a App) Stop(ctx context.Context) {
	a.logger.Info("Stopping app...")

	doneCh := make(chan error)
	go func() {
		doneCh <- a.server.Stop(ctx)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			a.logger.Error("Error while stopping server: %v", slog.Any("error_details", err))
		}

		a.logger.Info("App has been stopped gracefully")

	case <-ctx.Done():
		a.logger.Warn("App stopped forced")
	}
}
