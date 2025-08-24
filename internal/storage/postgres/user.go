package postgres

import (
	"context"
	"fmt"

	"github.com/Ostmind/tgbot/internal/helpers"
	"github.com/Ostmind/tgbot/internal/models"
)

func (store *Storage) GetUserByTelegramID(ctx context.Context, id string) (user models.User, err error) {
	sqlStatement := `SELECT * FROM public.users where telegram_id =$1;`

	rows, err := store.DB.Query(ctx, sqlStatement, id)
	if err != nil {
		return user, fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	err = rows.Scan(&user.ID, &user.TelegramID, &user.UserName, &user.Password, &user.Created)
	if err != nil {
		return user, fmt.Errorf("failed to parse DB %w", err)
	}

	return user, nil
}

func (store *Storage) DeleteUser(ctx context.Context, id string) error {
	sqlStatement := `DELETE FROM public.users WHERE id = $1;`

	result, err := store.DB.Exec(ctx, sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (store *Storage) AddUser(ctx context.Context,
	telegramID string,
	userName string,
	password string) (id string, refreshToken string, err error) {
	sqlStatement := `INSERT INTO public.users
					(telegram_id,username,created_at,password,refresh_token) 
					values ($1,$2,now(),$3,$4);`

	hashedPassword, err := helpers.ValidatePassword(password)
	if err != nil {
		return "", "", err
	}

	hashedRefreshToken, err := helpers.NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	result, err := store.DB.Exec(ctx, sqlStatement, telegramID, userName, hashedPassword, hashedRefreshToken)
	if err != nil {
		if !result.Insert() {
			return "", "", models.ErrUnique
		}

		return "", "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.categories where telegram_id = $1;`

	rows, err := store.DB.Query(ctx, sqlStatement, telegramID)
	if err != nil {
		return "", "", fmt.Errorf("failed to query DB %w", err)
	}

	defer rows.Close()

	err = rows.Scan(&id)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse DB %w", err)
	}

	return id, refreshToken, nil
}
