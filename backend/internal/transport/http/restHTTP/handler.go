package restHTTP

import (
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/auth"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/coefficients"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/service"
)

type Deps struct {
	Logger             *slog.Logger
	ProductService     product.ProductService
	CategoryService    category.CategoryService
	AuthService        auth.AuthService
	PresetService      preset.PresetService
	AttributeService   attribute.AttributeService
	CoefficientService coefficients.CoefficientService
	ServiceService     service.ServiceService
}

func NewDeps(
	Logger *slog.Logger,
	ProductService product.ProductService,
	CategoryService category.CategoryService,
	AuthService auth.AuthService,
	PresetService preset.PresetService,
	AttributeService attribute.AttributeService,
	CoefficientService coefficients.CoefficientService,
	ServiceService service.ServiceService,
) (*Deps, error) {
	if ProductService == nil {
		return nil, fmt.Errorf("missing ProductService dependency")
	}
	if CategoryService == nil {
		return nil, fmt.Errorf("missing CategoryService dependency")
	}
	if AuthService == nil {
		return nil, fmt.Errorf("missing AuthService dependency")
	}
	if PresetService == nil {
		return nil, fmt.Errorf("missing PresetService dependency")
	}
	if AttributeService == nil {
		return nil, fmt.Errorf("missing AttributeService dependency")
	}
	if CoefficientService == nil {
		return nil, fmt.Errorf("missing CoefficientService dependency")
	}
	if ServiceService == nil {
		return nil, fmt.Errorf("missing ServiceService dependency")
	}
	if Logger == nil {
		return nil, fmt.Errorf("missing Logger dependency")
	}
	return &Deps{
		Logger:             Logger,
		ProductService:     ProductService,
		CategoryService:    CategoryService,
		AuthService:        AuthService,
		PresetService:      PresetService,
		AttributeService:   AttributeService,
		CoefficientService: CoefficientService,
		ServiceService:     ServiceService,
	}, nil
}

type Handlers struct {
	ProductHandler      *product.Handler
	CategoryHandler     *category.Handler
	AuthHandler         *auth.Handler
	PresetHandler       *preset.Handler
	AttributeHandler    *attribute.Handler
	CoefficientsHandler *coefficients.Handler
	ServiceHandler      *service.Handler
}

func New(deps *Deps) (*Handlers, error) {
	if deps == nil {
		return nil, fmt.Errorf("missing deps for handlers")
	}

	// product handler
	prodDeps, err := product.NewDeps(deps.Logger, deps.ProductService)
	if err != nil {
		return nil, fmt.Errorf("product handler init: %w", err)
	}
	prodHandler := product.New(prodDeps)

	// category handler
	catDeps, err := category.NewDeps(deps.Logger, deps.CategoryService)
	if err != nil {
		return nil, fmt.Errorf("category handler init: %w", err)
	}
	catHandler := category.New(catDeps)

	// auth handler
	authDeps, err := auth.NewDeps(deps.Logger, deps.AuthService)
	if err != nil {
		return nil, fmt.Errorf("auth handler init: %w", err)
	}
	authHandler := auth.New(authDeps)

	// preset handler
	presetDeps, err := preset.NewDeps(deps.Logger, deps.PresetService)
	if err != nil {
		return nil, fmt.Errorf("preset handler init: %w", err)
	}
	presetHandler := preset.New(presetDeps)

	// attribute handler
	attrDeps, err := attribute.NewDeps(deps.Logger, deps.AttributeService)
	if err != nil {
		return nil, fmt.Errorf("attribute handler init: %w", err)
	}
	attrHandler := attribute.New(attrDeps)

	// coefficients handler
	coeffDeps, err := coefficients.NewDeps(deps.Logger, deps.CoefficientService)
	if err != nil {
		return nil, fmt.Errorf("coefficients handler init: %w", err)
	}
	coeffHandler := coefficients.New(coeffDeps)

	// service handler
	serviceDeps, err := service.NewDeps(deps.Logger, deps.ServiceService)
	if err != nil {
		return nil, fmt.Errorf("service handler init: %w", err)
	}
	serviceHandler := service.New(serviceDeps)

	return &Handlers{
		ProductHandler:      prodHandler,
		CategoryHandler:     catHandler,
		AuthHandler:         authHandler,
		PresetHandler:       presetHandler,
		AttributeHandler:    attrHandler,
		CoefficientsHandler: coeffHandler,
		ServiceHandler:      serviceHandler,
	}, nil
}
