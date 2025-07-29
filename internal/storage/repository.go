package storage

import (
	"context"
	"tgbot/internal/models"
)

type UserRepository interface {
	AddUser(ctx context.Context, telegramID string, userName string) (id string, err error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
