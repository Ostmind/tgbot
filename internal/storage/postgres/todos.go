package postgres

import (
	"context"
	"fmt"
)

type ToDoList struct {
	ID          string
	TelegramID  string
	Title       string
	Description string
	Completed   string
}

func (store *Storage) GetUserTodos(ctx context.Context, telegramID string) ([]ToDoList, error) {
	query := `SELECT id, telegram_id, title, description, completed FROM public.todos WHERE telegram_id = $1`

	rows, err := store.DB.Query(ctx, query, telegramID)
	if err != nil {
		return nil, fmt.Errorf("query user todos: %w", err)
	}
	defer rows.Close()

	var todos []ToDoList
	for rows.Next() {
		var t ToDoList
		if err := rows.Scan(&t.ID, &t.TelegramID, &t.Title, &t.Description, &t.Completed); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}
