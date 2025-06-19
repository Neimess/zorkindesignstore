package restHTTP

type CategoryService interface {
}

type CategoryHandler struct{}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{}
}
