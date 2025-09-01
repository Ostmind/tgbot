package todos

import (
	"errors"
	"github.com/Ostmind/tgbot/internal/storage"
	"github.com/Ostmind/tgbot/internal/todo/config"
	"github.com/Ostmind/tgbot/internal/todo/models"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type TodoController struct {
	Manager storage.TodoRepository
	cfg     *config.AppConfig
	logger  *slog.Logger
}

func NewTodoHandler(manager storage.TodoRepository, cfg *config.AppConfig, log *slog.Logger) *TodoController {
	return &TodoController{manager, cfg, log}
}

func (ctr TodoController) AddTodo(echo echo.Context) error {
	ctr.logger.Debug("Post Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	desc := echo.Param("desc")

	id, err := ctr.Manager.AddTodo(echo.Request().Context(), telegramID, title, desc)
	if err != nil {
		if errors.Is(err, models.ErrUnique) {
			return echo.NoContent(http.StatusConflict)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.JSON(http.StatusOK, id)
}

func (ctr TodoController) DeleteTodo(echo echo.Context) error {
	ctr.logger.Debug("Delete Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	err := ctr.Manager.DeleteTodo(echo.Request().Context(), telegramID, title)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}

func (ctr TodoController) UpdateTodo(echo echo.Context) error {
	ctr.logger.Debug("Update Request for todos")

	telegramID := echo.Param("telegramID")

	title := echo.Param("title")

	isDone := echo.Param("isDone")

	isDoneBool := false

	if isDone == "true" {
		isDoneBool = true
	}

	err := ctr.Manager.UpdateTodo(echo.Request().Context(), telegramID, title, isDoneBool)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return echo.NoContent(http.StatusNotFound)
		}

		return echo.NoContent(http.StatusInternalServerError)
	}

	return echo.NoContent(http.StatusOK)
}
