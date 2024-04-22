package testutil

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	Dir string
	DB  *sql.DB
}

func NewMigrator(dir string, db *sql.DB) *Migrator {
	return &Migrator{
		Dir: dir,
		DB:  db,
	}
}

func (m Migrator) MigrateUp() error {
	migrations := &migrate.FileMigrationSource{
		Dir: m.Dir,
	}

	_, err := migrate.Exec(m.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}

func (m Migrator) MigrateDown() error {
	migrations := &migrate.FileMigrationSource{
		Dir: m.Dir,
	}

	_, err := migrate.Exec(m.DB, "mysql", migrations, migrate.Down)
	if err != nil {
		return err
	}

	return nil
}
