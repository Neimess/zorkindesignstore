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
)

func newHandler(mockSvc *mocks.MockCategoryService) *Handler {
	dep, err := NewDeps(slog.New(slog.DiscardHandler), mockSvc)
	if err != nil {
		return nil
	}
	h := New(dep)
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
	h := newHandler(mockSvc)

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

func TestCreate_Errors(t *testing.T) {
	cases := []struct {
		name           string
		body           string
		svcErr         error
		wantCode       int
		wantMsg        string
		wantValidation bool
	}{
		{
			name:     "invalid JSON",
			body:     `{`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid JSON",
		},
		{
			name:           "validation error - empty name",
			body:           `{"name":"","description":"test"}`,
			wantCode:       http.StatusUnprocessableEntity,
			wantValidation: true,
		},
		{
			name:     "service error",
			body:     `{"name":"test","description":"test"}`,
			svcErr:   errors.New("service error"),
			wantCode: http.StatusInternalServerError,
			wantMsg:  "internal server error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockCategoryService(t)
			if tc.svcErr != nil {
				mockSvc.
					EXPECT().
					CreateCategory(mock.Anything, mock.Anything).
					Return(nil, tc.svcErr).
					Once()
			}

			h := newHandler(mockSvc)
			req := httptest.NewRequest(http.MethodPost, "/api/admin/category", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			h.CreateCategory(w, req)

			assert.Equal(t, tc.wantCode, w.Code)

			if tc.wantValidation {
				var resp httputils.ValidationErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Errors)
			} else {
				var resp httputils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Contains(t, resp.Message, tc.wantMsg)
			}
		})
	}
}

func TestGetCategory_Success(t *testing.T) {
	mockSvc := mocks.NewMockCategoryService(t)
	h := newHandler(mockSvc)

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
			h := newHandler(mockSvc)

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
	h := newHandler(mockSvc)

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
			h := newHandler(mockSvc)

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

func TestUpdate_Errors(t *testing.T) {
	cases := []struct {
		name           string
		body           string
		vars           map[string]string
		svcErr         error
		wantCode       int
		wantMsg        string
		wantValidation bool
	}{
		{
			name:     "invalid JSON",
			body:     `{`,
			vars:     map[string]string{"id": "1"},
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid JSON",
		},
		{
			name:     "invalid category id",
			body:     `{"name":"test","description":"test"}`,
			vars:     map[string]string{"id": "abc"},
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid category id",
		},
		{
			name:           "validation error - empty name",
			body:           `{"name":"","description":"test"}`,
			vars:           map[string]string{"id": "1"},
			wantCode:       http.StatusUnprocessableEntity,
			wantValidation: true,
		},
		{
			name:     "service error",
			body:     `{"name":"test","description":"test"}`,
			vars:     map[string]string{"id": "1"},
			svcErr:   errors.New("service error"),
			wantCode: http.StatusInternalServerError,
			wantMsg:  "internal server error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockCategoryService(t)
			if tc.svcErr != nil {
				mockSvc.
					EXPECT().
					UpdateCategory(mock.Anything, mock.Anything).
					Return(nil, tc.svcErr).
					Once()
			}

			h := newHandler(mockSvc)
			req := httptest.NewRequest(http.MethodPut, "/api/admin/category/1", strings.NewReader(tc.body))
			req = withChiParams(req, tc.vars)
			w := httptest.NewRecorder()

			h.UpdateCategory(w, req)

			assert.Equal(t, tc.wantCode, w.Code)

			if tc.wantValidation {
				var resp httputils.ValidationErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Errors)
			} else {
				var resp httputils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Contains(t, resp.Message, tc.wantMsg)
			}
		})
	}
}
