package model

type TestItemResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InsertTestItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListTestItemsRequest struct {
	PageSize int `json:"pageSize"`
	Offset   int `json:"offset"`
}
