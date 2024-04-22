package store

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //nolint:revive
	"github.com/jmoiron/sqlx"

	"github.com/acernik/go-server-boilerplate/internal/app"
	svcmodel "github.com/acernik/go-server-boilerplate/internal/service/model"
)

type Store interface {
	InsertTestTableItem(ttItem svcmodel.TestTableItem) (int64, error)
}

type store struct {
	db         *sqlx.DB
	insertStmt *sqlx.NamedStmt
}

func NewStore(cfg *app.Config) (Store, error) {
	dbDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sqlx.Connect("mysql", dbDSN)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	str := &store{db: db}

	insertStmt, err := str.PrepareInsertStatement()
	if err != nil {
		return nil, err
	}
	str.insertStmt = insertStmt

	return str, nil
}

func (s *store) InsertTestTableItem(ttItem svcmodel.TestTableItem) (int64, error) {
	queryParams := map[string]interface{}{ //nolint:gofmt
		"name":        ttItem.Name,
		"description": ttItem.Description,
	}

	res, err := s.insertStmt.Exec(queryParams)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
