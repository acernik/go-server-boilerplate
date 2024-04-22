package mapper

import (
	"github.com/acernik/go-server-boilerplate/internal/server/model"
	svcmodel "github.com/acernik/go-server-boilerplate/internal/service/model"
)

func MapTestTableItemAPIToService(req model.InsertTestItemRequest) svcmodel.TestTableItem {
	return svcmodel.TestTableItem{
		Name:        req.Name,
		Description: req.Description,
	}
}
