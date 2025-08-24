package todos

import (
	"context"
	"fmt"
	"github.com/Ostmind/tgbot/internal/models"
	"github.com/Ostmind/tgbot/internal/storage"
)

type StorageToDo struct {
	storage storage.ToDoRepository
}

func New(storage storage.ToDoRepository) *StorageToDo {
	return &StorageToDo{
		storage: storage,
	}
}

func (c StorageToDo) GetUserTodos(ctx context.Context, telegramID string) (res []models.ToDoList, err error) {
	res, err = c.storage.GetUserTodos(ctx, telegramID)
	if err != nil {
		return res, fmt.Errorf("failed to add user %w", err)
	}

	return res, nil
}

func (c StorageToDo) AddToDo(ctx context.Context, telegramID string, title string, desc string) (id string, err error) {
	id, err = c.storage.AddToDo(ctx, telegramID, title, desc)
	if err != nil {
		return id, fmt.Errorf("failed to get user %w", err)
	}

	return id, nil
}

func (c StorageToDo) DeleteToDo(ctx context.Context, telegramID string, title string) error {
	err := c.storage.DeleteToDo(ctx, telegramID, title)
	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil
}

func (c StorageToDo) UpdateToDo(ctx context.Context, telegramID string, title string, isDone bool) error {
	err := c.storage.UpdateToDo(ctx, telegramID, title, isDone)
	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil
}
