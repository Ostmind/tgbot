package postgres

import (
	"context"
	"fmt"

	"github.com/Ostmind/tgbot/internal/models"
)

func (store *Storage) GetUserTodos(ctx context.Context, telegramID string) ([]models.ToDoList, error) {
	query := `SELECT id, telegram_id, title, description, completed FROM public.todos WHERE telegram_id = $1`

	rows, err := store.DB.Query(ctx, query, telegramID)
	if err != nil {
		return nil, fmt.Errorf("query user todos: %w", err)
	}
	defer rows.Close()

	var todos []models.ToDoList
	for rows.Next() {
		var t models.ToDoList
		if err := rows.Scan(&t.ID, &t.TelegramID, &t.Title, &t.Description, &t.Completed); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (store *Storage) AddToDo(ctx context.Context, telegramID string, title string, desc string) (id string, err error) {
	sqlStatement := `INSERT INTO public.todos
					(telegram_id,title,description) 
					values ($1,$2,$3);`

	result, err := store.DB.Exec(ctx, sqlStatement, telegramID, title, desc)
	if err != nil {
		if !result.Insert() {
			return "", models.ErrUnique
		}

		return "", fmt.Errorf("error adding to DB %w", err)
	}

	sqlStatement = `SELECT id FROM public.todos where telegram_id = $1 and title = $2;`

	rows, err := store.DB.Query(ctx, sqlStatement, telegramID, title)
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

func (store *Storage) DeleteToDo(ctx context.Context, telegramID string, title string) error {
	sqlStatement := `DELETE FROM public.todos WHERE telegram_id = $1 and title = $2;`

	result, err := store.DB.Exec(ctx, sqlStatement, telegramID, title)
	if err != nil {
		return fmt.Errorf("error deleting from DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (store *Storage) UpdateToDo(ctx context.Context, telegramID string, title string, isDone bool) error {
	sqlStatement := `UPDATE public.todos SET completed=$1 WHERE telegram_id = $2 and title = $3;`

	result, err := store.DB.Exec(ctx, sqlStatement, isDone, telegramID, title)
	if err != nil {
		return fmt.Errorf("error updating DB %w", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
