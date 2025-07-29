package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tgbot/internal/models"

	"github.com/labstack/echo/v4"
)

type UserManager interface {
	AddUser(ctx context.Context, telegramID string, userName string) (id string, err error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserController struct {
	manager UserManager
	logger  *slog.Logger
}

func NewUserHandler(manager UserManager, log *slog.Logger) *UserController {
	return &UserController{manager, log}
}

func (ctr UserController) GetUserByID(echo echo.Context) error {
	ctr.logger.Debug("Get Request for User")

	userID := echo.Param("userId")

	res, err := ctr.manager.GetUserByID(echo.Request().Context(), userID)
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr UserController) AddUser(echo echo.Context) error {
	ctr.logger.Debug("Post Request for User")

	telegramID := echo.Param("telegramID")

	userName := echo.Param("userName")

	res, err := ctr.manager.AddUser(echo.Request().Context(), telegramID, userName)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, res)
}

func (ctr UserController) DeleteUser(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for User")

	userID := echo.Param("userId")

	err := ctr.manager.DeleteUser(echo.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
