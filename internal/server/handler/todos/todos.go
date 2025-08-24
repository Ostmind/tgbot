package todos

import (
	"errors"
	"github.com/Ostmind/tgbot/internal/config"
	"github.com/Ostmind/tgbot/internal/models"
	"github.com/Ostmind/tgbot/internal/storage"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type ToDoController struct {
	Manager storage.ToDoRepository
	cfg     *config.AppConfig
	logger  *slog.Logger
}

func NewToDoHandler(manager storage.ToDoRepository, cfg *config.AppConfig, log *slog.Logger) *ToDoController {
	return &ToDoController{manager, cfg, log}
}

func (ctr ToDoController) AddToDo(echo echo.Context) error {
	ctr.logger.Debug("Post Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	desc := echo.Param("desc")

	id, err := ctr.Manager.AddToDo(echo.Request().Context(), telegramID, title, desc)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, id)
}

func (ctr ToDoController) DeleteToDo(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	err := ctr.Manager.DeleteToDo(echo.Request().Context(), telegramID, title)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}

func (ctr ToDoController) UpdateToDo(echo echo.Context) error {
	ctr.logger.Debug("Update Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	isDone := echo.Param("isDone")

	isDoneBool := false

	if isDone == "true" {
		isDoneBool = true
	}

	err := ctr.Manager.UpdateToDo(echo.Request().Context(), telegramID, title, isDoneBool)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
