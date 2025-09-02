package todos

import (
	"context"
	"fmt"
	"github.com/Ostmind/tgbot/internal/storage"
	"github.com/Ostmind/tgbot/internal/todo/models"
)

type StorageTodo struct {
	storage storage.TodoRepository
}

func New(storage storage.TodoRepository) *StorageTodo {
	return &StorageTodo{
		storage: storage,
	}
}

func (c StorageTodo) GetUserTodos(ctx context.Context, telegramID string) (res []models.ToDoList, err error) {
	res, err = c.storage.GetUserTodos(ctx, telegramID)
	if err != nil {
		return res, fmt.Errorf("failed to add user %w", err)
	}

	return res, nil
}

func (c StorageTodo) AddTodo(ctx context.Context, telegramID string, title string, desc string) (id string, err error) {
	id, err = c.storage.AddTodo(ctx, telegramID, title, desc)
	if err != nil {
		return id, fmt.Errorf("failed to get user %w", err)
	}

	return id, nil
}

func (c StorageTodo) DeleteTodo(ctx context.Context, telegramID string, title string) error {
	err := c.storage.DeleteTodo(ctx, telegramID, title)
	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil
}

func (c StorageTodo) UpdateTodo(ctx context.Context, telegramID string, title string, isDone bool) error {
	err := c.storage.UpdateTodo(ctx, telegramID, title, isDone)
	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil
}
