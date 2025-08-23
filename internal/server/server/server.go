package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/server/handler/user"
	"github.com/Ostmind/tgbot/internal/server/middleware"
	"github.com/Ostmind/tgbot/internal/storage/postgres"

	"github.com/labstack/echo/v4"
)

type Server struct {
	server  *echo.Echo
	logger  *slog.Logger
	storage *postgres.Storage
	port    int
}

func New(logger *slog.Logger,
	cfg *config.AppConfig,
	db *postgres.Storage,
	userHandler *user.Controller) *Server {
	server := echo.New()

	server.Use(middleware.LogRequestAndAuthenticateUser(logger, userHandler, cfg.Auth))

	categoryGroup := server.Group("users")

	categoryGroup.GET("", userHandler.GetUserByTelegramID)
	categoryGroup.DELETE("/:userId", userHandler.DeleteUser)
	categoryGroup.POST("/create/:telegramID", userHandler.AddUser)

	return &Server{
		logger:  logger,
		server:  server,
		storage: db,
		port:    cfg.Srv.Port,
	}
}
func (s Server) Run() {
	s.logger.Info("Server is running on: localhost", "Port", s.port)

	if err := s.server.Start(fmt.Sprintf("localhost:%d", s.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server starting error: %v", slog.Any("error_details", err))
		}
	}
}

func (s Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping DB Connection")

	s.storage.Close()

	s.logger.Info("Stopping server...")
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error("Error: ", slog.Any("error_details", err))

		return fmt.Errorf("error while stopping Server Request %w", err)
	}

	return nil
}
