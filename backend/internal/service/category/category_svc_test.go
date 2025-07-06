package category_test

import (
	"context"
	"errors"
	"testing"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	catservice "github.com/Neimess/zorkin-store-project/internal/service/category"
	"github.com/Neimess/zorkin-store-project/internal/service/category/mocks"
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
	type testCase struct {
		name      string
		input     *category.Category
		mockSetup func()
		expectErr bool
		expectNil bool
	}

	tests := []testCase{
		{
			name:  "success",
			input: &category.Category{ID: 1, Name: "Test"},
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, &category.Category{ID: 1, Name: "Test"}).Return(&category.Category{ID: 1, Name: "Test"}, nil).Once()
			},
			expectErr: false,
			expectNil: false,
		},
		{
			name:      "validation error",
			input:     &category.Category{ID: 1, Name: ""},
			mockSetup: func() {},
			expectErr: true,
			expectNil: true,
		},
		{
			name:  "repository error",
			input: &category.Category{ID: 2, Name: "RepoFail"},
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, &category.Category{ID: 2, Name: "RepoFail"}).Return(nil, errors.New("db error")).Once()
			},
			expectErr: true,
			expectNil: true,
		},
		{
			name:  "double creation",
			input: &category.Category{ID: 3, Name: "Double"},
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, &category.Category{ID: 3, Name: "Double"}).Return(&category.Category{ID: 3, Name: "Double"}, nil).Twice()
			},
			expectErr: false,
			expectNil: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest() // reset mocks for each subtest
			tc.mockSetup()
			if tc.name == "double creation" {
				created1, err1 := s.svc.CreateCategory(context.Background(), tc.input)
				created2, err2 := s.svc.CreateCategory(context.Background(), tc.input)
				s.NoError(err1)
				s.NoError(err2)
				s.Equal(tc.input, created1)
				s.Equal(tc.input, created2)
			} else {
				created, err := s.svc.CreateCategory(context.Background(), tc.input)
				if tc.expectErr {
					s.Error(err)
				} else {
					s.NoError(err)
				}
				if tc.expectNil {
					s.Nil(created)
				} else {
					s.Equal(tc.input, created)
				}
			}
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CategoryServiceSuite) TestGetCategory() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expect    *category.Category
		expectErr bool
	}

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&category.Category{ID: 1, Name: "Test"}, nil).Once()
			},
			expect:    &category.Category{ID: 1, Name: "Test"},
			expectErr: false,
		},
		{
			name: "not found",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("GetByID", mock.Anything, int64(2)).Return(nil, category.ErrCategoryNotFound).Once()
			},
			expect:    nil,
			expectErr: true,
		},
		{
			name: "repository error",
			id:   3,
			mockSetup: func() {
				s.mockRepo.On("GetByID", mock.Anything, int64(3)).Return(nil, errors.New("db error")).Once()
			},
			expect:    nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			got, err := s.svc.GetCategory(context.Background(), tc.id)
			if tc.expectErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tc.expect, got)
			}
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CategoryServiceSuite) TestUpdateCategory() {
	type testCase struct {
		name      string
		input     *category.Category
		mockSetup func()
		expect    *category.Category
		expectErr bool
	}

	tests := []testCase{
		{
			name:  "success",
			input: &category.Category{ID: 1, Name: "Updated"},
			mockSetup: func() {
				inputted := category.Category{ID: 1, Name: "Updated"}
				s.mockRepo.On("Update", mock.Anything, &inputted).Return(&category.Category{ID: 1, Name: "Updated"}, nil).Once()
			},
			expect:    &category.Category{ID: 1, Name: "Updated"},
			expectErr: false,
		},
		{
			name:      "validation error",
			input:     &category.Category{ID: 1, Name: ""},
			mockSetup: func() {},
			expect:    nil,
			expectErr: true,
		},
		{
			name:  "not found",
			input: &category.Category{ID: 2, Name: "Updated"},
			mockSetup: func() {
				inputted := category.Category{ID: 2, Name: "Updated"}
				s.mockRepo.On("Update", mock.Anything, &inputted).Return(nil, category.ErrCategoryNotFound).Once()
			},
			expect:    nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			updated, err := s.svc.UpdateCategory(context.Background(), tc.input)
			if tc.expectErr {
				s.Error(err)
				s.Nil(updated)
			} else {
				s.NoError(err)
				s.Equal(tc.expect, updated)
			}
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CategoryServiceSuite) TestDeleteCategory() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expectErr bool
	}

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name: "not found",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(2)).Return(category.ErrCategoryNotFound).Once()
			},
			expectErr: false,
		},
		{
			name: "repository error",
			id:   3,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(3)).Return(errors.New("db error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			err := s.svc.DeleteCategory(context.Background(), tc.id)
			if tc.expectErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *CategoryServiceSuite) TestListCategories() {
	type testCase struct {
		name      string
		mockSetup func()
		expect    []category.Category
		expectErr bool
	}

	tests := []testCase{
		{
			name: "success",
			mockSetup: func() {
				s.mockRepo.On("List", mock.Anything).Return([]category.Category{{ID: 1, Name: "Test"}}, nil).Once()
			},
			expect:    []category.Category{{ID: 1, Name: "Test"}},
			expectErr: false,
		},
		{
			name: "empty list",
			mockSetup: func() {
				s.mockRepo.On("List", mock.Anything).Return([]category.Category{}, nil).Once()
			},
			expect:    []category.Category{},
			expectErr: false,
		},
		{
			name: "repository error",
			mockSetup: func() {
				s.mockRepo.On("List", mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expect:    nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			got, err := s.svc.ListCategories(context.Background())
			if tc.expectErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tc.expect, got)
			}
			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}
