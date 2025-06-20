// internal/adapter/http/resthandler/deps.go
package restHTTP

import "log/slog"

type Deps struct {
	ProductService  ProductService
	CategoryService CategoryService
	AuthService     AuthService
	Logger          *slog.Logger
}
