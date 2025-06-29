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

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/mocks"
)

type fakeValidator struct{ err error }

func (f fakeValidator) StructCtx(ctx context.Context, s interface{}) error {
	return f.err
}

func newHandler(mockSvc *mocks.MockProductService, valErr error) *Handler {
	dep, err := NewDeps(validator.New(), slog.New(slog.DiscardHandler), mockSvc)
	require.NoError(&testing.T{}, err)
	h := New(dep)
	h.val = fakeValidator{err: valErr}
	return h
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func TestCreate_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	reqBody := dto.ProductCreateRequest{Name: "Product X"}
	raw, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/product", bytes.NewReader(raw))
	w := httptest.NewRecorder()

	mockSvc.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(int64(1), nil).
		Once()

	h.Create(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreate_Errors(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		valErr   error
		svcErr   error
		wantCode int
	}{
		{"bad JSON", `{`, nil, nil, http.StatusBadRequest},
		{"validation failed", `{"name":"P"}`, errors.New("validation"), nil, http.StatusUnprocessableEntity},
		{"bad category", `{"name":"P"}`, nil, prodDom.ErrBadCategoryID, http.StatusBadRequest},
		{"invalid attr", `{"name":"P"}`, nil, prodDom.ErrInvalidAttribute, http.StatusUnprocessableEntity},
		{"generic error", `{"name":"P"}`, nil, errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockProductService(t)
			h := newHandler(mockSvc, tc.valErr)

			req := httptest.NewRequest(http.MethodPost, "/api/admin/product", bytes.NewReader([]byte(tc.body)))
			w := httptest.NewRecorder()

			if tc.valErr == nil && tc.svcErr != nil && tc.body[0] == '{' {
				mockSvc.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(int64(0), tc.svcErr).
					Once()
			}

			h.Create(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}

func TestCreateWithAttributes_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	reqBody := dto.ProductCreateRequest{Name: "Product Y"}
	raw, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/product/detailed", bytes.NewReader(raw))
	w := httptest.NewRecorder()

	mockSvc.EXPECT().
		CreateWithAttrs(mock.Anything, mock.Anything).
		Return(int64(1), nil).
		Once()

	h.CreateWithAttributes(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateWithAttributes_Errors(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	// Bad JSON
	req := httptest.NewRequest(http.MethodPost, "/api/admin/product/detailed", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	h.CreateWithAttributes(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Service returns invalid attribute
	reqBody := dto.ProductCreateRequest{Name: "Product Z"}
	raw, _ := json.Marshal(reqBody)
	req = httptest.NewRequest(http.MethodPost, "/api/admin/product/detailed", bytes.NewReader(raw))
	w = httptest.NewRecorder()

	mockSvc.EXPECT().
		CreateWithAttrs(mock.Anything, mock.Anything).
		Return(1, prodDom.ErrInvalidAttribute).
		Once()

	h.CreateWithAttributes(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDetailed_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/product/1", nil)
	req = withChiParams(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	mockSvc.EXPECT().
		GetDetailed(mock.Anything, int64(1)).
		Return(&prodDom.Product{ID: 1, Name: "Product 1"}, nil).
		Once()

	h.GetDetailed(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDetailed_Errors(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	// Bad ID
	req := httptest.NewRequest(http.MethodGet, "/api/product/abc", nil)
	req = withChiParams(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()
	h.GetDetailed(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Not found
	req = httptest.NewRequest(http.MethodGet, "/api/product/1", nil)
	req = withChiParams(req, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	mockSvc.EXPECT().
		GetDetailed(mock.Anything, int64(1)).
		Return(nil, prodDom.ErrProductNotFound).
		Once()
	h.GetDetailed(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListByCategory_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/product/category/2", nil)
	req = withChiParams(req, map[string]string{"id": "2"})
	w := httptest.NewRecorder()

	mockSvc.EXPECT().
		GetByCategoryID(mock.Anything, int64(2)).
		Return([]prodDom.Product{{ID: 2, Name: "ProdCat"}}, nil).
		Once()

	h.ListByCategory(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	body := dto.ProductUpdateRequest{Name: "Updated"}
	raw, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/product/7", bytes.NewReader(raw))
	req = withChiParams(req, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	mockSvc.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()
	h.Update(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDelete_Success(t *testing.T) {
	mockSvc := mocks.NewMockProductService(t)
	h := newHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/product/5", nil)
	req = withChiParams(req, map[string]string{"id": "5"})
	w := httptest.NewRecorder()

	mockSvc.EXPECT().Delete(mock.Anything, int64(5)).Return(nil).Once()
	h.Delete(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
