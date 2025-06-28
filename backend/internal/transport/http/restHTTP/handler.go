package restHTTP

import (
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/auth"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product"
	"github.com/go-playground/validator/v10"
)

type Deps struct {
	Logger           *slog.Logger
	Validator        *validator.Validate
	ProductService   product.ProductService
	CategoryService  category.CategoryService
	AuthService      auth.AuthService
	PresetService    preset.PresetService
	AttributeService attribute.AttributeService
}

type Handlers struct {
	ProductHandler   *product.Handler
	CategoryHandler  *category.Handler
	AuthHandler      *auth.Handler
	PresetHandler    *preset.Handler
	AttributeHandler *attribute.Handler
}

func New(deps *Deps) (*Handlers, error) {
	if deps == nil {
		return nil, fmt.Errorf("missing deps for handlers")
	}
	// ensure logger
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	// product handler
	prodDeps, err := product.NewDeps(deps.ProductService, deps.Validator)
	if err != nil {
		return nil, fmt.Errorf("product handler init: %w", err)
	}
	prodHandler := product.New(prodDeps)

	// category handler
	catDeps, err := category.NewDeps(deps.CategoryService)
	if err != nil {
		return nil, fmt.Errorf("category handler init: %w", err)
	}
	catHandler := category.New(catDeps)

	// auth handler
	authDeps, err := auth.NewDeps(deps.AuthService)
	if err != nil {
		return nil, fmt.Errorf("auth handler init: %w", err)
	}
	authHandler := auth.New(authDeps)

	// preset handler
	presetDeps, err := preset.NewDeps(deps.PresetService)
	if err != nil {
		return nil, fmt.Errorf("preset handler init: %w", err)
	}
	presetHandler := preset.New(presetDeps)

	// attribute handler
	attrDeps, err := attribute.NewDeps(deps.AttributeService)
	if err != nil {
		return nil, fmt.Errorf("attribute handler init: %w", err)
	}
	attrHandler := attribute.New(attrDeps)

	return &Handlers{
		ProductHandler:   prodHandler,
		CategoryHandler:  catHandler,
		AuthHandler:      authHandler,
		PresetHandler:    presetHandler,
		AttributeHandler: attrHandler,
	}, nil
}
