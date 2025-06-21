package restHTTP

import "log/slog"

type Deps struct {
	Logger                   *slog.Logger
	ProductService           ProductService
	CategoryService          CategoryService
	AuthService              AuthService
	CategoryAttributeService CategoryAttributeService
	PresetService            PresetService
}

type Handlers struct {
	ProductHandler           *ProductHandler
	CategoryHandler          *CategoryHandler
	AuthHandler              *AuthHandler
	CategoryAttributeHandler *CategoryAttributeHandler
	PresetHandler            *PresetHandler
}

func New(deps *Deps) *Handlers {
	return &Handlers{
		ProductHandler:           NewProductHandler(deps.ProductService, deps.Logger),
		CategoryHandler:          NewCategoryHandler(deps.CategoryService, deps.Logger),
		AuthHandler:              NewAuthHandler(deps.AuthService, deps.Logger),
		CategoryAttributeHandler: NewCategoryAttributeHandler(deps.CategoryAttributeService, deps.Logger),
		PresetHandler:            NewPresetHandler(deps.PresetService, deps.Logger),
	}
}
