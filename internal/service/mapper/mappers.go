package mapper

import (
	srvmodel "github.com/acernik/go-server-boilerplate/internal/service/model"
	"github.com/acernik/go-server-boilerplate/internal/store/model"
)

func MapTestTableItemsDBToDomain(dbResult ...model.TestTableItem) []srvmodel.TestTableItem {
	domainResult := make([]srvmodel.TestTableItem, len(dbResult))
	for i, item := range dbResult {
		domainResult[i] = srvmodel.TestTableItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		}
	}

	return domainResult
}
