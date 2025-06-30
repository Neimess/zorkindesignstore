package category_test

import (
	"context"
	"testing"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	catservice "github.com/Neimess/zorkin-store-project/internal/service/category"
	"github.com/Neimess/zorkin-store-project/internal/service/category/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CategoryServiceSuite struct {
	suite.Suite
	svc      *catservice.Service
	mockRepo *mocks.MockCategoryRepository
	logger   *slog.Logger
}

func (s *CategoryServiceSuite) SetupTest() {
	s.mockRepo = new(mocks.MockCategoryRepository)
	s.logger = slog.Default()
	deps, _ := catservice.NewDeps(s.mockRepo, s.logger)
	s.svc = catservice.New(deps)
}

func (s *CategoryServiceSuite) TestCreateCategory() {
	s.Run("успешное создание", func() {
		cat := &category.Category{ID: 1, Name: "Test"}
		s.mockRepo.On("Create", mock.Anything, cat).Return(cat, nil).Once()

		created, err := s.svc.CreateCategory(context.Background(), cat)
		s.NoError(err)
		s.Equal(cat, created)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("ошибка валидации", func() {
		cat := &category.Category{ID: 1, Name: ""}
		created, err := s.svc.CreateCategory(context.Background(), cat)
		s.Error(err)
		s.Nil(created)
	})
}

func (s *CategoryServiceSuite) TestGetCategory() {
	s.Run("успешное получение", func() {
		cat := &category.Category{ID: 1, Name: "Test"}
		s.mockRepo.On("GetByID", mock.Anything, int64(1)).Return(cat, nil).Once()

		got, err := s.svc.GetCategory(context.Background(), 1)
		s.NoError(err)
		s.Equal(cat, got)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("категория не найдена", func() {
		s.mockRepo.On("GetByID", mock.Anything, int64(2)).Return(nil, category.ErrCategoryNotFound).Once()

		got, err := s.svc.GetCategory(context.Background(), 2)
		s.ErrorIs(err, category.ErrCategoryNotFound)
		s.Nil(got)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *CategoryServiceSuite) TestUpdateCategory() {
	s.Run("успешное обновление", func() {
		cat := &category.Category{ID: 1, Name: "Updated"}
		s.mockRepo.On("Update", mock.Anything, int64(1), "Updated").Return(nil).Once()

		err := s.svc.UpdateCategory(context.Background(), cat)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("ошибка валидации", func() {
		cat := &category.Category{ID: 1, Name: ""}
		err := s.svc.UpdateCategory(context.Background(), cat)
		s.Error(err)
	})

	s.Run("категория не найдена", func() {
		cat := &category.Category{ID: 2, Name: "Updated"}
		s.mockRepo.On("Update", mock.Anything, int64(2), "Updated").Return(category.ErrCategoryNotFound).Once()

		err := s.svc.UpdateCategory(context.Background(), cat)
		s.ErrorIs(err, category.ErrCategoryNotFound)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *CategoryServiceSuite) TestDeleteCategory() {
	s.Run("успешное удаление", func() {
		s.mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		err := s.svc.DeleteCategory(context.Background(), 1)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("категория не найдена", func() {
		s.mockRepo.On("Delete", mock.Anything, int64(2)).Return(category.ErrCategoryNotFound).Once()

		err := s.svc.DeleteCategory(context.Background(), 2)
		s.ErrorIs(err, category.ErrCategoryNotFound)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *CategoryServiceSuite) TestListCategories() {
	s.Run("успешный список", func() {
		cats := []category.Category{{ID: 1, Name: "Test"}}
		s.mockRepo.On("List", mock.Anything).Return(cats, nil).Once()

		got, err := s.svc.ListCategories(context.Background())
		s.NoError(err)
		s.Equal(cats, got)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("ошибка репозитория", func() {
		s.mockRepo.On("List", mock.Anything).Return(nil, assert.AnError).Once()

		got, err := s.svc.ListCategories(context.Background())
		s.Error(err)
		s.Nil(got)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}
