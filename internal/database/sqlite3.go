package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func New(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_busy_timeout=5000", dbPath))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	slog.Info("Database connection successful", "path", dbPath)

	slog.Info("Initializing database...")
	err = initialize(db, "todos", Todo{})
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		panic(err)
	}
	slog.Info("Database initialized!")

	slog.Info("Database connection successful", "path", dbPath)
	return db
}

func initialize(db *sql.DB, tableName string, model interface{}) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct")
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		columnType := ""

		slog.Debug("Field", "name", field.Name, "type", field.Type.Kind(), "tag", columnName)
		if columnName == "" {
			continue // Skip fields that don't have a `db` tag
		}

		switch field.Type.Kind() {
		case reflect.Int:
			if field.Name == "ID" {
				columnType = "INTEGER PRIMARY KEY AUTOINCREMENT"
			} else {
				columnType = "INTEGER"
			}
		case reflect.String:
			columnType = "TEXT"
		case reflect.Bool:
			columnType = "BOOLEAN"
		case reflect.Struct:
			columnType = "TIMESTAMP"
		default:
			return fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}

		query += fmt.Sprintf("%s %s,", columnName, columnType)
	}

	// Remove trailing comma and close the statement
	query = query[:len(query)-1] + ");"
	slog.Debug("Query", "query", query)

	// Execute the statement
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
