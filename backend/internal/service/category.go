package service

import "log/slog"

type CategoryRepository interface {
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{}
}
