package storage

import (
	"context"
	"tgbot/internal/models"
)

type Repository interface {
	AddUser(ctx context.Context, telegramID string, userName string, password string) (id string, refreshToken string, err error)
	GetUserByTelegramID(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
