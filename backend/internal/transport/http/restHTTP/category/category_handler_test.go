package category

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"
)

type fakeValidator struct{ err error }

func (f fakeValidator) StructCtx(ctx context.Context, s interface{}) error {
	return f.err
}

func newHandler(mockSvc *mocks.MockCategoryService, valErr error) *Handler {
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

func TestCreateCategory_Success(t *testing.T) {
	mockSvc := mocks.NewMockCategoryService(t)
	h := newHandler(mockSvc, nil)

	reqBody := map[string]interface{}{"name": "Books"}
	raw, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/category", bytes.NewReader(raw))
	w := httptest.NewRecorder()

	mockSvc.
		EXPECT().
		CreateCategory(mock.Anything, mock.MatchedBy(func(c *catDom.Category) bool {
			return c.Name == "Books"
		})).
		Return(&catDom.Category{ID: 42, Name: "Books"}, nil).
		Once()

	h.CreateCategory(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/api/category/42", w.Header().Get("Location"))
	var resp dto.CategoryResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(42), resp.ID)
	assert.Equal(t, "Books", resp.Name)
	mockSvc.AssertExpectations(t)
}

func TestCreateCategory_Errors(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		valErr   error
		svcErr   error
		wantCode int
		wantMsg  string
	}{
		{
			name:     "invalid JSON",
			body:     `{`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid JSON",
		},
		{
			name:     "validation failure",
			body:     `{"name":""}`,
			valErr:   errors.New("any"),
			wantCode: http.StatusUnprocessableEntity,
			wantMsg:  "invalid product data",
		},
		{
			name:     "empty name",
			body:     `{"name":" "}`,
			svcErr:   catDom.ErrCategoryNameEmpty,
			wantCode: http.StatusUnprocessableEntity,
			wantMsg:  "category name cannot be empty",
		},
		{
			name:     "name too long",
			body:     `{"name":"` + strings.Repeat("x", 300) + `"}`,
			svcErr:   catDom.ErrCategoryNameTooLong,
			wantCode: http.StatusUnprocessableEntity,
			wantMsg:  "category name must be at most 255 characters",
		},
		{
			name:     "generic service error",
			body:     `{"name":"Books"}`,
			svcErr:   errors.New("boom"),
			wantCode: http.StatusInternalServerError,
			wantMsg:  "internal server error",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockSvc := mocks.NewMockCategoryService(t)
			h := newHandler(mockSvc, tc.valErr)

			req := httptest.NewRequest(http.MethodPost, "/api/admin/category", bytes.NewReader([]byte(tc.body)))
			w := httptest.NewRecorder()

			// only stub service if validation passed and JSON was valid
			if tc.valErr == nil && tc.svcErr != nil && tc.body[0] == '{' {
				mockSvc.
					EXPECT().
					CreateCategory(mock.Anything, mock.Anything).
					Return(nil, tc.svcErr).
					Once()
			}

			h.CreateCategory(w, req)
			assert.Equal(t, tc.wantCode, w.Code)

			// on error we return JSON with { "message": "<msg>" }
			var resp httputils.ErrorResponse
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tc.wantMsg, resp.Message)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestGetCategory_Success(t *testing.T) {
	mockSvc := mocks.NewMockCategoryService(t)
	h := newHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/category/7", nil)
	req = withChiParams(req, map[string]string{"id": "7"})
	w := httptest.NewRecorder()

	mockSvc.
		EXPECT().
		GetCategory(mock.Anything, int64(7)).
		Return(&catDom.Category{ID: 7, Name: "C7"}, nil).
		Once()

	h.GetCategory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var out dto.CategoryResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.Equal(t, int64(7), out.ID)
	assert.Equal(t, "C7", out.Name)
}

func TestGetCategory_Errors(t *testing.T) {
	cases := []struct {
		name     string
		param    string
		svcErr   error
		wantCode int
		wantMsg  string
	}{
		{"bad id", "abc", nil, http.StatusBadRequest, "invalid id"},
		{"not found", "9", catDom.ErrCategoryNotFound, http.StatusNotFound, "category not found"},
		{"internal error", "5", errors.New("boom"), http.StatusInternalServerError, "internal server error"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockCategoryService(t)
			h := newHandler(mockSvc, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/category/"+tc.param, nil)
			req = withChiParams(req, map[string]string{"id": tc.param})
			w := httptest.NewRecorder()

			if tc.svcErr != nil && tc.param != "abc" {
				mockSvc.
					EXPECT().
					GetCategory(mock.Anything, mock.Anything).
					Return(nil, tc.svcErr).
					Once()
			}

			h.GetCategory(w, req)
			assert.Equal(t, tc.wantCode, w.Code)

			var resp httputils.ErrorResponse
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tc.wantMsg, resp.Message)
		})
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	mockSvc := mocks.NewMockCategoryService(t)
	h := newHandler(mockSvc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/category/3", nil)
	req = withChiParams(req, map[string]string{"id": "3"})
	w := httptest.NewRecorder()

	mockSvc.
		EXPECT().
		DeleteCategory(mock.Anything, int64(3)).
		Return(nil).
		Once()

	h.DeleteCategory(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteCategory_Errors(t *testing.T) {
	cases := []struct {
		name     string
		param    string
		svcErr   error
		wantCode int
		wantMsg  string
	}{
		{"bad id", "0", nil, http.StatusBadRequest, "invalid category id"},
		{"not found", "4", catDom.ErrCategoryNotFound, http.StatusNotFound, "category not found"},
		{"in use", "5", catDom.ErrCategoryInUse, http.StatusConflict, "category is in use and cannot be deleted"},
		{"other error", "6", errors.New("boom"), http.StatusInternalServerError, "internal server error"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockCategoryService(t)
			h := newHandler(mockSvc, nil)

			req := httptest.NewRequest(http.MethodDelete, "/api/admin/category/"+tc.param, nil)
			req = withChiParams(req, map[string]string{"id": tc.param})
			w := httptest.NewRecorder()

			if tc.svcErr != nil && tc.param != "0" {
				mockSvc.
					EXPECT().
					DeleteCategory(mock.Anything, mock.Anything).
					Return(tc.svcErr).
					Once()
			}

			h.DeleteCategory(w, req)
			assert.Equal(t, tc.wantCode, w.Code)

			var resp httputils.ErrorResponse
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			// for 204 there is no body
			if tc.wantCode != http.StatusNoContent {
				assert.Equal(t, tc.wantMsg, resp.Message)
			}
		})
	}
}
