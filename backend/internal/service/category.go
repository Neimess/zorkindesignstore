package service

import "log/slog"

type CategoryRepository interface {
}

type CategoryService struct {
	repo CategoryRepository
	log  *slog.Logger
}

func NewCategoryService(repo CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{}
}
