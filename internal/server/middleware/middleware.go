package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/server/auth"
	"github.com/Ostmind/tgbot/internal/server/handler/user"

	"github.com/labstack/echo/v4"
)

func LogRequestAndAuthenticateUser(logger *slog.Logger, manager *user.Controller, cfgAuth config.AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echo echo.Context) error {
			start := time.Now()

			if !(echo.Request().Method == http.MethodPost && strings.Contains(echo.Request().URL.String(), "/users/create")) {
				err := auth.Authentication(echo, manager, cfgAuth)
				if err != nil {
					logger.Error("Request error",
						"error", err,
						"method", echo.Request().Method,
						"url", echo.Request().URL.String())

					echo.Error(err)

					return nil
				}
			}

			err := next(echo)

			stop := time.Now()

			logger.Info("Request: ",
				"Method", echo.Request().Method,
				"URL", echo.Request().URL,
				"Time", stop.Sub(start),
				"Http Code", echo.Response().Status)

			return err
		}
	}
}
