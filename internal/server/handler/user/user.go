package user

import (
	"errors"
	"log/slog"
	"net/http"
	"tgbot/internal/models"
	"tgbot/internal/storage"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Manager storage.Repository
	logger  *slog.Logger
}

func NewUserHandler(manager storage.Repository, log *slog.Logger) *Controller {
	return &Controller{manager, log}
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

	userName := echo.Param("userName")

	password := echo.Param("password")

	res, err := ctr.Manager.AddUser(echo.Request().Context(), telegramID, userName, password)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

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
