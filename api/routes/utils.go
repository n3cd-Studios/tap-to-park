package routes

type PaginatorOutput[T interface{}] struct {
	Items []T   `json:"items"`
	Pages int64 `json:"pages"`
	Page  int64 `json:"page"`
}
