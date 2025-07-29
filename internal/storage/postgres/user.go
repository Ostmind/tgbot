package postgres

import (
	"context"
	"fmt"
	"tgbot/internal/models"
)

type Users struct {
	db *Storage
}

func NewUsers(db *Storage) (*Users, error) {
	return &Users{
		db: db,
	}, nil
}

func (c *Users) GetUserByID(ctx context.Context, id string) (user models.User, err error) {
	sqlStatement := `SELECT * FROM public.users where id =$1`

	rows, err := c.db.DB.Query(ctx, sqlStatement, id)
	if err != nil {
		return user, fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	err = rows.Scan(&user.ID, &user.TelegramID, &user.UserName, &user.AuthSalt, &user.Created)

	if err != nil {
		return user, fmt.Errorf("failed to parse DB %w", err)
	}

	return user, nil
}

func (c *Users) DeleteUser(ctx context.Context, id string) error {
	sqlStatement := `DELETE FROM public.users WHERE id = $1;`

	result, err := c.db.DB.Exec(ctx, sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (c *Users) AddUser(ctx context.Context, telegramID string, userName string) (id string, err error) {
	sqlStatement := `INSERT INTO public.users
					(telegram_id,username,created_at) 
					values ($1,$2,now());`

	result, err := c.db.DB.Exec(ctx, sqlStatement, telegramID, userName)

	if err != nil {
		if !result.Insert() {
			return "", models.ErrUnique
		}

		return "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.categories where telegram_id = $1`

	rows, err := c.db.DB.Query(ctx, sqlStatement, telegramID)
	if err != nil {
		return "", fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	err = rows.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to parse DB %w", err)
	}

	return id, nil
}
