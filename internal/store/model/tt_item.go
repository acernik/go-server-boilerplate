package model

type TestTableItem struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
