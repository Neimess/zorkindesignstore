package product

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

type ProductHandlerSuite struct {
	suite.Suite
	h       *Handler
	mockSvc *mocks.MockProductService
}

func (s *ProductHandlerSuite) SetupTest() {
	s.mockSvc = mocks.NewMockProductService(s.T())
	dep, err := NewDeps(slog.New(slog.NewTextHandler(io.Discard, nil)), s.mockSvc)
	s.Require().NoError(err)
	s.h = New(dep)
}

func (s *ProductHandlerSuite) TestCreate() {
	type testCase struct {
		name      string
		body      interface{}
		svcMock   func(*mocks.MockProductService)
		wantCode  int
		wantCheck func(*testing.T, *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name: "success",
			body: dto.ProductRequest{
				Name:       "ValidProduct",
				Price:      10,
				CategoryID: 1,
				Attributes: []dto.ProductAttributeRequest{{Name: "Объем", Value: "10"}},
			},
			svcMock: func(s *mocks.MockProductService) {
				s.EXPECT().Create(mock.Anything, mock.Anything).Return(&prodDom.Product{ID: 1, Name: "ValidProduct", Price: 10, CategoryID: 1}, nil).Once()
			},
			wantCode: http.StatusCreated,
			wantCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp dto.ProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), resp.ProductID)
				assert.Equal(t, "ValidProduct", resp.Name)
			},
		},
		{
			name:     "bad JSON",
			body:     `{`,
			svcMock:  nil,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation error - empty name",
			body: dto.ProductRequest{
				Name:       "",
				Price:      10,
				CategoryID: 1,
			},
			svcMock:  nil,
			wantCode: http.StatusUnprocessableEntity,
			wantCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp httputils.ValidationErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Errors)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()

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
			if tc.svcMock != nil {
				s.mockSvc.AssertExpectations(s.T())
			}
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
			id:   "2",
			svcMock: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetDetailed(mock.Anything, int64(2)).Return(nil, prodDom.ErrProductNotFound).Once()
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
				req = httptest.NewRequest(http.MethodGet, "/api/product/"+tc.id, nil)
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
		{
			name:     "bad id",
			id:       "abc",
			svcMock:  nil,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			var req *http.Request
			if tc.name == "bad id" {
				req = httptest.NewRequest(http.MethodGet, "/api/product/category/abc", nil)
				req = withChiParams(req, map[string]string{"id": "abc"})
			} else {
				req = httptest.NewRequest(http.MethodGet, "/api/product/category/"+tc.id, nil)
				req = withChiParams(req, map[string]string{"id": tc.id})
			}
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
	body := dto.ProductRequest{Name: "UpdatedProduct", Price: 10, CategoryID: 1}
	raw, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/product/7", bytes.NewReader(raw))
	req = withChiParams(req, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	s.SetupTest()
	s.mockSvc.EXPECT().Update(mock.Anything, mock.Anything).Return(&prodDom.Product{ID: 7, Name: "UpdatedProduct", Price: 10, CategoryID: 1}, nil).Once()
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
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}
