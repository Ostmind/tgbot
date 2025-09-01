package user

import (
	"errors"
	"github.com/Ostmind/tgbot/internal/todo/config"
	helpers2 "github.com/Ostmind/tgbot/internal/todo/helpers"
	"github.com/Ostmind/tgbot/internal/todo/models"
	"log/slog"
	"net/http"

	"github.com/Ostmind/tgbot/internal/storage"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Manager storage.Repository
	cfg     *config.AppConfig
	logger  *slog.Logger
}

func NewUserHandler(manager storage.Repository, cfg *config.AppConfig, log *slog.Logger) *Controller {
	return &Controller{manager, cfg, log}
}

func (ctr Controller) GetUserByTelegramID(echo echo.Context) error {
	ctr.logger.Debug("Get Request for User")

	userID := echo.Param("userId")

	res, err := ctr.Manager.GetUserByTelegramID(echo.Request().Context(), userID)
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr Controller) AddUser(echo echo.Context) error {
	ctr.logger.Debug("Post Request for User")

	telegramID := echo.Param("telegramID")

	userNameCookie, _ := echo.Cookie("username")

	passwordCookie, _ := echo.Cookie("password")

	res, refreshToken, err := ctr.Manager.AddUser(echo.Request().Context(), telegramID, userNameCookie.Value, passwordCookie.Value)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	jwt, err := helpers2.GenerateJWT(userNameCookie.Value, ctr.cfg.Auth.JWTAccessTokenTTL, ctr.cfg.Auth.JWTSecret)
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	cookie := helpers2.SetCookie("AccessToken", jwt, ctr.cfg.Auth.JWTAccessTokenTTL, false)

	echo.SetCookie(cookie)

	cookie = helpers2.SetCookie("RefreshToken", refreshToken, ctr.cfg.Auth.JWTRefreshTokenTTL, true)

	echo.SetCookie(cookie)

	return echo.JSON(http.StatusOK, res)
}

func (ctr Controller) DeleteUser(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for User")

	userID := echo.Param("userId")

	err := ctr.Manager.DeleteUser(echo.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
