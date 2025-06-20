package restHTTP

type Handlers struct {
	ProductHandler  *ProductHandler
	CategoryHandler *CategoryHandler
	AuthHandler     *AuthHandler
}

func New(deps *Deps) *Handlers {
	return &Handlers{
		ProductHandler:  NewProductHandler(deps.ProductService, deps.Logger),
		CategoryHandler: NewCategoryHandler(deps.CategoryService, deps.Logger),
		AuthHandler:     NewAuthHandler(deps.AuthService, deps.Logger),
	}
}
