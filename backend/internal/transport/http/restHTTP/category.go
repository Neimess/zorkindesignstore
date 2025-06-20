package restHTTP

import "log/slog"

type CategoryService interface {
}

type CategoryHandler struct {
	log *slog.Logger
	srv CategoryService
}

func NewCategoryHandler(service CategoryService, log *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		srv: service,
		log: log,
	}
}
