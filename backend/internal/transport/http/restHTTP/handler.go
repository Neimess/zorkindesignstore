package restHTTP

import (
	"log/slog"

	ch "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
)

type Deps struct {
	Logger          *slog.Logger
	ProductService  ProductService
	CategoryService ch.CategoryService
	AuthService     AuthService
	PresetService   PresetService
}

type Handlers struct {
	ProductHandler  *ProductHandler
	CategoryHandler *ch.CategoryHandler
	AuthHandler     *AuthHandler
	PresetHandler   *PresetHandler
}

func New(deps *Deps) *Handlers {
	return &Handlers{
		ProductHandler:  NewProductHandler(deps.ProductService, deps.Logger),
		CategoryHandler: ch.NewCategoryHandler(deps.CategoryService),
		AuthHandler:     NewAuthHandler(deps.AuthService, deps.Logger),
		PresetHandler:   NewPresetHandler(deps.PresetService, deps.Logger),
	}
}
