package middleware

import (
	"log/slog"
	"tgbot/internal/server/handler/user"
	"tgbot/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func LogRequest(logger *slog.Logger, manager *user.Controller) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echo echo.Context) error {
			start := time.Now()

			userName := echo.Param("userName")

			password := echo.Param("password")

			userRepository, err := manager.Manager.GetUserByTelegramID(echo.Request().Context(), userName)
			if err != nil {
				logger.Error("Request error",
					"error", err,
					"method", echo.Request().Method,
					"url", echo.Request().URL.String())

				echo.Error(err)

				return nil
			}

			err = utils.ComparePassword(userRepository.Password, password)
			if err != nil {
				logger.Error("Request error",
					"error", err,
					"method", echo.Request().Method,
					"url", echo.Request().URL.String())

				echo.Error(err)

				return nil
			}

			err = next(echo)

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
