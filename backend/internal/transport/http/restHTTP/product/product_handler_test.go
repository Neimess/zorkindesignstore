package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/mocks"
)

type fakeValidator struct{ err error }

func (f fakeValidator) StructCtx(ctx context.Context, s interface{}) error {
	return f.err
}

type ProductHandlerSuite struct {
	suite.Suite
	h       *Handler
	mockSvc *mocks.MockProductService
	valErr  error
}

func (s *ProductHandlerSuite) SetupTest() {
	s.mockSvc = mocks.NewMockProductService(s.T())
	s.valErr = nil
	dep, err := NewDeps(validator.New(), slog.New(slog.DiscardHandler), s.mockSvc)
	s.Require().NoError(err)
	s.h = New(dep)
	s.h.val = fakeValidator{err: s.valErr}
}

func (s *ProductHandlerSuite) TestCreate() {
	type testCase struct {
		name      string
		body      interface{}
		valErr    error
		svcMock   func(*mocks.MockProductService)
		wantCode  int
		wantCheck func(*testing.T, *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name:   "success",
			body:   dto.ProductCreateRequest{Name: "Product X", Price: 78.3, CategoryID: 1},
			valErr: nil,
			svcMock: func(s *mocks.MockProductService) {
				s.EXPECT().Create(mock.Anything, mock.Anything).Return(&prodDom.Product{ID: 1, Name: "test"}, nil).Once()
			},
			wantCode: http.StatusCreated,
			wantCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dto.ProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, int64(1), resp.ProductID)
				require.Equal(t, "test", resp.Name)
			},
		},
		{
			name:      "bad JSON",
			body:      `{`,
			valErr:    nil,
			svcMock:   nil,
			wantCode:  http.StatusBadRequest,
			wantCheck: nil,
		},
		{
			name:      "validation failed",
			body:      dto.ProductCreateRequest{Name: ""},
			valErr:    errors.New("validation"),
			svcMock:   nil,
			wantCode:  http.StatusUnprocessableEntity,
			wantCheck: nil,
		},
		{
			name:   "bad category",
			body:   dto.ProductCreateRequest{Name: "ValidName", Price: 10, CategoryID: 999},
			valErr: nil,
			svcMock: func(s *mocks.MockProductService) {
				s.EXPECT().Create(mock.Anything, mock.Anything).Return(&prodDom.Product{}, prodDom.ErrBadCategoryID).Once()
			},
			wantCode:  http.StatusBadRequest,
			wantCheck: nil,
		},
		{
			name:   "invalid attr",
			body:   dto.ProductCreateRequest{Name: "ValidName", Price: 10, CategoryID: 1},
			valErr: nil,
			svcMock: func(s *mocks.MockProductService) {
				s.EXPECT().Create(mock.Anything, mock.Anything).Return(&prodDom.Product{}, prodDom.ErrInvalidAttribute).Once()
			},
			wantCode:  http.StatusUnprocessableEntity,
			wantCheck: nil,
		},
		{
			name:   "generic error",
			body:   dto.ProductCreateRequest{Name: "ValidName", Price: 10, CategoryID: 1},
			valErr: nil,
			svcMock: func(s *mocks.MockProductService) {
				s.EXPECT().Create(mock.Anything, mock.Anything).Return(&prodDom.Product{}, errors.New("boom")).Once()
			},
			wantCode:  http.StatusInternalServerError,
			wantCheck: nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			s.valErr = tc.valErr
			s.h.val = fakeValidator{err: tc.valErr}

			var req *http.Request
			if s, ok := tc.body.(string); ok {
				req = httptest.NewRequest(http.MethodPost, "/api/admin/product", bytes.NewReader([]byte(s)))
			} else {
				b, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(http.MethodPost, "/api/admin/product", bytes.NewReader(b))
			}
			w := httptest.NewRecorder()

			if tc.svcMock != nil {
				tc.svcMock(s.mockSvc)
			}

			s.h.Create(w, req)
			assert.Equal(s.T(), tc.wantCode, w.Code)
			if tc.wantCheck != nil {
				tc.wantCheck(s.T(), w)
			}
		})
	}
}

func (s *ProductHandlerSuite) TestCreateWithAttributes() {
	type testCase struct {
		name     string
		body     interface{}
		svcMock  func(*mocks.MockProductService)
		wantCode int
	}

	tests := []testCase{
		{
			name: "success",
			body: dto.ProductCreateRequest{Name: "Product Y"},
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().CreateWithAttrs(mock.Anything, mock.Anything).Return(&prodDom.Product{ID: 1}, nil).Once()
			},
			wantCode: http.StatusCreated,
		},
		{
			name:     "bad JSON",
			body:     `{`,
			svcMock:  nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid attribute",
			body: dto.ProductCreateRequest{Name: "Product Z"},
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().CreateWithAttrs(mock.Anything, mock.Anything).Return(&prodDom.Product{ID: 1}, prodDom.ErrInvalidAttribute).Once()
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			var req *http.Request
			if s, ok := tc.body.(string); ok {
				req = httptest.NewRequest(http.MethodPost, "/api/admin/product/detailed", bytes.NewReader([]byte(s)))
			} else {
				b, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(http.MethodPost, "/api/admin/product/detailed", bytes.NewReader(b))
			}
			w := httptest.NewRecorder()
			if tc.svcMock != nil {
				tc.svcMock(s.mockSvc)
			}
			s.h.CreateWithAttributes(w, req)
			assert.Equal(s.T(), tc.wantCode, w.Code)
		})
	}
}

func (s *ProductHandlerSuite) TestGetDetailed() {
	type testCase struct {
		name     string
		id       string
		svcMock  func(*mocks.MockProductService)
		wantCode int
	}

	tests := []testCase{
		{
			name: "success",
			id:   "1",
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetDetailed(mock.Anything, int64(1)).Return(&prodDom.Product{ID: 1, Name: "Product 1"}, nil).Once()
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "bad id",
			id:       "abc",
			svcMock:  nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "1",
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetDetailed(mock.Anything, int64(1)).Return(nil, prodDom.ErrProductNotFound).Once()
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			var req *http.Request
			if tc.name == "bad id" {
				req = httptest.NewRequest(http.MethodGet, "/api/product/abc", nil)
				req = withChiParams(req, map[string]string{"id": "abc"})
			} else {
				req = httptest.NewRequest(http.MethodGet, "/api/product/1", nil)
				req = withChiParams(req, map[string]string{"id": tc.id})
			}
			w := httptest.NewRecorder()
			if tc.svcMock != nil {
				tc.svcMock(s.mockSvc)
			}
			s.h.GetDetailed(w, req)
			assert.Equal(s.T(), tc.wantCode, w.Code)
		})
	}
}

func (s *ProductHandlerSuite) TestListByCategory() {
	type testCase struct {
		name     string
		id       string
		svcMock  func(*mocks.MockProductService)
		wantCode int
	}

	tests := []testCase{
		{
			name: "success",
			id:   "2",
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetByCategoryID(mock.Anything, int64(2)).Return([]prodDom.Product{{ID: 2, Name: "ProdCat"}}, nil).Once()
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			req := httptest.NewRequest(http.MethodGet, "/api/product/category/"+tc.id, nil)
			req = withChiParams(req, map[string]string{"id": tc.id})
			w := httptest.NewRecorder()
			if tc.svcMock != nil {
				tc.svcMock(s.mockSvc)
			}
			s.h.ListByCategory(w, req)
			assert.Equal(s.T(), tc.wantCode, w.Code)
		})
	}
}

func (s *ProductHandlerSuite) TestUpdate() {
	body := dto.ProductUpdateRequest{Name: "Updated"}
	raw, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/product/7", bytes.NewReader(raw))
	req = withChiParams(req, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	s.SetupTest()
	s.mockSvc.EXPECT().Update(mock.Anything, mock.Anything).Return(validProduct(), nil).Once()
	s.h.Update(w, req)
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *ProductHandlerSuite) TestDelete() {
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/product/5", nil)
	req = withChiParams(req, map[string]string{"id": "5"})
	w := httptest.NewRecorder()

	s.SetupTest()
	s.mockSvc.EXPECT().Delete(mock.Anything, int64(5)).Return(nil).Once()
	s.h.Delete(w, req)
	assert.Equal(s.T(), http.StatusNoContent, w.Code)
}

func TestProductHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerSuite))
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func validProduct() *prodDom.Product {
	desc := "desc"
	img := "img"
	return &prodDom.Product{
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
