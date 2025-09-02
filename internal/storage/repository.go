package storage

import (
	"context"
	"github.com/Ostmind/tgbot/internal/todo/models"
)

type Repository interface {
	AddUser(ctx context.Context, telegramID string, userName string, password string) (id string, refreshToken string, err error)
	GetUserByTelegramID(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type TodoRepository interface {
	AddTodo(ctx context.Context, telegramID string, title string, desc string) (id string, err error)
	DeleteTodo(ctx context.Context, telegramID string, title string) error
	UpdateTodo(ctx context.Context, telegramID string, title string, isDone bool) error
	GetUserTodos(ctx context.Context, telegramID string) ([]models.ToDoList, error)
}
