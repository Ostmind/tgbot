package app

import (
	"context"
	"fmt"
	"log/slog"
	"tgbot/internal/config"
	userhandler "tgbot/internal/server/handler/user"
	srv "tgbot/internal/server/server"
	"tgbot/internal/services/users"
	"tgbot/internal/storage"
	"tgbot/internal/storage/postgres"
	"time"
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

	categoryStorage, err := postgres.NewUsers(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't create categories %w", err)
	}

	userManager := users.New(categoryStorage)

	userHandler := userhandler.NewUserHandler(userManager, logger)

	server := srv.New(logger, &cfg.Srv, db, userHandler)

	return &App{
		server: server,
		logger: logger,
		db:     db,
		cfg:    cfg,
	}, nil
}

func (a App) Run() {
	a.logger.Info("Starting app...")
	err := storage.RunMigration(a.db, a.logger, a.cfg.Srv.MigrationPath)

	if err != nil {
		a.logger.Error("couldn't run migrations %w", slog.Any("error_details", err))
	}

	a.server.Run()
}

func (a App) Stop(ctx context.Context, shutdownTimeout time.Duration) {
	a.logger.Info("Stopping app...")

	timeout := shutdownTimeout

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	doneCh := make(chan error)
	go func() {
		doneCh <- a.server.Stop(ctxWithTimeout)
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
