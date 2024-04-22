package service

import (
	"github.com/acernik/go-server-boilerplate/internal/service/model"
	str "github.com/acernik/go-server-boilerplate/internal/store"
)

type Service interface {
	InsertTestTableItem(ttItem model.TestTableItem) (int64, error)
}

type service struct {
	store str.Store
}

func NewService(store str.Store) Service {
	return &service{store: store}
}

func (s *service) InsertTestTableItem(ttItem model.TestTableItem) (int64, error) {
	return s.store.InsertTestTableItem(ttItem)
}
