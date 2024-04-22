package store

import "github.com/jmoiron/sqlx"

func (s *store) PrepareInsertStatement() (*sqlx.NamedStmt, error) {
	return s.db.PrepareNamed(`
		INSERT INTO test_table (name, description) VALUES (:name, :description);
	`)
}
