package users

import (
	"context"
	"fmt"
	"tgbot/internal/models"
	"tgbot/internal/storage"
)

type StorageUser struct {
	storage storage.Repository
}

func New(storage storage.Repository) *StorageUser {
	return &StorageUser{
		storage: storage,
	}
}

func (c StorageUser) AddUser(ctx context.Context, telegramID string, userName string, password string) (id string, refreshToken string, err error) {
	id, refreshToken, err = c.storage.AddUser(ctx, telegramID, userName, password)

	if err != nil {
		return id, "", fmt.Errorf("failed to add user %w", err)
	}

	return id, refreshToken, nil
}

func (c StorageUser) GetUserByTelegramID(ctx context.Context, id string) (models.User, error) {
	user, err := c.storage.GetUserByTelegramID(ctx, id)

	if err != nil {
		return user, fmt.Errorf("failed to get user %w", err)
	}

	return user, nil
}

func (c StorageUser) DeleteUser(ctx context.Context, id string) error {
	err := c.storage.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %w", err)
	}

	return nil
}
