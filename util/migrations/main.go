package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

const DOWN = "down"
const UP = "up"

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonHandler)

	migrations := &migrate.FileMigrationSource{
		Dir: "database/migrations/api",
	}

	var migrationDirection migrate.MigrationDirection

	arg := os.Args[1]

	switch arg {
	case DOWN:
		migrationDirection = migrate.Down
	case UP:
		migrationDirection = migrate.Up
	default:
		logger.Warn("not supported migration")
		return
	}

	dsn := os.Args[2]
	if len(dsn) == 0 {
		logger.Error("DB_DSN environment variable not set")
		return
	}

	db, _ := sql.Open("mysql", dsn)
	_, err := migrate.Exec(db, "mysql", migrations, migrationDirection)
	if err != nil {
		logger.Error("migration failed with error: %s", err)
		return
	}

	logger.Info("successfully applied migrations")
}
