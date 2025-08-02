package models

import "time"

type User struct {
	ID         string    `db:"id"`
	TelegramID string    `db:"telegram_id"`
	UserName   string    `db:"username"`
	Password   string    `db:"password"`
	Created    time.Time `db:"created_at"`
}
