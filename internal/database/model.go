package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type Todo struct {
	ID          int          `db:"id"`
	Description string       `db:"description"`
	Completed   bool         `db:"completed"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	DB          *sql.DB
}

func NewTodo(db *sql.DB) *Todo {
	return &Todo{
		DB: db,
	}
}

func (t Todo) Create(description string, completed bool) error {
	atTheMoment := time.Now().Format("2006-01-02 15:04:05")
	_, err := t.DB.Exec("INSERT INTO todos (description, completed, created_at) VALUES (?, ?, ?)", description, completed, atTheMoment)
	if err != nil {
		return fmt.Errorf("failed to save todo: %w", err)
	}
	return nil
}

func (t Todo) Delete(id int) error {
	_, err := t.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

func (t Todo) Update(id int, completed bool) error {
	atTheMoment := time.Now().Format("2006-01-02 15:04:05")
	_, err := t.DB.Exec("UPDATE todos SET completed = ?, updated_at = ? WHERE id = ?", completed, atTheMoment, id)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

// TODO: Implement Find method
// This method should find a todo by using finder with '/' as the keypress to activate the function
func (t Todo) Find() error {
	return nil
}

func (t Todo) All() []Todo {
	var result []Todo
	todos, err := t.DB.Query("SELECT * FROM todos")
	if err != nil {
		return nil
	}
	//parse result to []Todo

	for todos.Next() {
		var todo Todo
		err := todos.Scan(&todo.ID, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt, &todo.DeletedAt)
		if err != nil {
			slog.Error("Error scanning todos", "error", err)
		}
		result = append(result, todo)
	}
	return result
}
