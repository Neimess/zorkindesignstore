package restHTTP

type Handlers struct {
	ProductHandler  *ProductHandler
	CategoryHandler *CategoryHandler
}

func New(deps *Deps) *Handlers {
	return &Handlers{
		ProductHandler:  NewProductHandler(deps.ProductService, deps.Logger),
		// CategoryHandler: NewCategoryHandler(deps.CategoryService, deps.Logger),
	}
}
