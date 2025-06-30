package product_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"log/slog"

	catdomain "github.com/Neimess/zorkin-store-project/internal/domain/category"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	productservice "github.com/Neimess/zorkin-store-project/internal/service/product"
	"github.com/Neimess/zorkin-store-project/internal/service/product/mocks"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductServiceSuite struct {
	suite.Suite
	svc      *productservice.Service
	mockRepo *mocks.MockProductRepository
	logger   *slog.Logger
}

func (s *ProductServiceSuite) SetupTest() {
	s.mockRepo = new(mocks.MockProductRepository)
	s.logger = slog.New(slog.DiscardHandler)
	deps, _ := productservice.NewDeps(s.mockRepo, s.logger)
	s.svc = productservice.New(deps)
}

func validProduct() *domProduct.Product {
	desc := "desc"
	img := "img"
	return &domProduct.Product{
		ID:          1,
		Name:        "Test",
		Price:       100.0,
		Description: &desc,
		CategoryID:  1,
		ImageURL:    &img,
		CreatedAt:   time.Now(),
		Attributes:  nil,
	}
}

func (s *ProductServiceSuite) TestCreate() {
	type testCase struct {
		name      string
		input     *domProduct.Product
		mockSetup func()
		expectErr error
	}

	tests := []testCase{
		{
			name:  "success",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(validProduct(), nil).Once()
			},
			expectErr: nil,
		},
		{
			name:  "bad category",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "bad category (category.ErrCategoryNotFound)",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, catdomain.ErrCategoryNotFound).Once()
			},
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "invalid attribute",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrValidation).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "invalid attribute (bad request)",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrBadRequest).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "repo error",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			prod, err := s.svc.Create(context.Background(), tc.input)
			if tc.expectErr == nil {
				s.NoError(err)
				s.NotNil(prod)
				s.Equal(int64(1), prod.ID)
			} else {
				s.Error(err)
				s.Nil(prod)
			}
		})
	}
}

func (s *ProductServiceSuite) TestCreateWithAttrs() {
	type testCase struct {
		name      string
		input     *domProduct.Product
		mockSetup func()
		expectErr error
	}

	tests := []testCase{
		{
			name:  "success",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(&domProduct.Product{ID: 1}, nil).Once()
			},
			expectErr: nil,
		},
		{
			name:  "bad category",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "bad category (category.ErrCategoryNotFound)",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, catdomain.ErrCategoryNotFound).Once()
			},
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "invalid attribute",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrValidation).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "invalid attribute (bad request)",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrBadRequest).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "repo error",
			input: validProduct(),
			mockSetup: func() {
				s.mockRepo.On("CreateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			prod, err := s.svc.CreateWithAttrs(context.Background(), tc.input)
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(int64(1), prod.ID)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *ProductServiceSuite) TestGetDetailed() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expectErr error
	}

	p := validProduct()

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("GetWithAttrs", mock.Anything, int64(1)).Return(p, nil).Once()
			},
			expectErr: nil,
		},
		{
			name: "not found",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("GetWithAttrs", mock.Anything, int64(2)).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: domProduct.ErrProductNotFound,
		},
		{
			name: "repo error",
			id:   3,
			mockSetup: func() {
				s.mockRepo.On("GetWithAttrs", mock.Anything, int64(3)).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			res, err := s.svc.GetDetailed(context.Background(), tc.id)
			if tc.expectErr == nil {
				s.NoError(err)
				s.NotNil(res)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *ProductServiceSuite) TestGetByCategoryID() {
	type testCase struct {
		name      string
		catID     int64
		mockSetup func()
		expect    []domProduct.Product
		expectErr error
	}

	p := validProduct()

	tests := []testCase{
		{
			name:  "success",
			catID: 1,
			mockSetup: func() {
				s.mockRepo.On("ListByCategory", mock.Anything, int64(1)).Return([]domProduct.Product{*p}, nil).Once()
			},
			expect:    []domProduct.Product{*p},
			expectErr: nil,
		},
		{
			name:  "not found",
			catID: 2,
			mockSetup: func() {
				s.mockRepo.On("ListByCategory", mock.Anything, int64(2)).Return(nil, der.ErrNotFound).Once()
			},
			expect:    nil,
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "not found (category.ErrCategoryNotFound)",
			catID: 3,
			mockSetup: func() {
				s.mockRepo.On("ListByCategory", mock.Anything, int64(3)).Return(nil, catdomain.ErrCategoryNotFound).Once()
			},
			expect:    nil,
			expectErr: domProduct.ErrBadCategoryID,
		},
		{
			name:  "repo error",
			catID: 4,
			mockSetup: func() {
				s.mockRepo.On("ListByCategory", mock.Anything, int64(4)).Return(nil, errors.New("db fail")).Once()
			},
			expect:    nil,
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			res, err := s.svc.GetByCategoryID(context.Background(), tc.catID)
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expect, res)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *ProductServiceSuite) TestUpdate() {
	type testCase struct {
		name      string
		input     *domProduct.Product
		mockSetup func()
		expectErr error
	}

	p := validProduct()

	tests := []testCase{
		{
			name:  "success",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(p, nil).Once()
			},
			expectErr: nil,
		},
		{
			name:  "not found",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: domProduct.ErrProductNotFound,
		},
		{
			name:  "not found (product.ErrProductNotFound)",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, domProduct.ErrProductNotFound).Once()
			},
			expectErr: domProduct.ErrProductNotFound,
		},
		{
			name:  "invalid attribute",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrValidation).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "invalid attribute (bad request)",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, der.ErrBadRequest).Once()
			},
			expectErr: domProduct.ErrInvalidAttribute,
		},
		{
			name:  "repo error",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("UpdateWithAttrs", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			prodRes, err := s.svc.Update(context.Background(), tc.input)
			if tc.expectErr == nil {
				s.Equal(tc.input.Name, prodRes.Name)
				s.NoError(err)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *ProductServiceSuite) TestDelete() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expectErr error
	}

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectErr: nil,
		},
		{
			name: "not found",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(2)).Return(der.ErrNotFound).Once()
			},
			expectErr: domProduct.ErrProductNotFound,
		},
		{
			name: "not found (product.ErrProductNotFound)",
			id:   3,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(3)).Return(domProduct.ErrProductNotFound).Once()
			},
			expectErr: domProduct.ErrProductNotFound,
		},
		{
			name: "repo error",
			id:   4,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(4)).Return(errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			err := s.svc.Delete(context.Background(), tc.id)
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
			}
		})
	}
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}
