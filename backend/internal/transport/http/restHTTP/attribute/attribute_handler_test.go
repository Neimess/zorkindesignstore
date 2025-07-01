package attribute

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

func newHandler(mockSvc *mocks.MockAttributeService) *Handler {
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

func TestCreateAttributesBatch_Success(t *testing.T) {
	mockSvc := mocks.NewMockAttributeService(t)
	h := newHandler(mockSvc)

	input := []map[string]interface{}{{"name": "A", "unit": "u"}}
	body, err := json.Marshal(map[string]interface{}{"data": input})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewReader(body))
	req = withChiParams(req, map[string]string{"categoryID": "1"})
	w := httptest.NewRecorder()

	mockSvc.
		On("CreateAttributesBatch", mock.Anything, int64(1), mock.Anything).
		Return(nil).
		Once()

	h.CreateAttributesBatch(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateAttributesBatch_Errors(t *testing.T) {
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
			vars:     map[string]string{"categoryID": "1"},
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid JSON",
		},
		{
			name:     "invalid category id",
			body:     `{"data":[{}]}`,
			vars:     map[string]string{"categoryID": "abc"},
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid category id",
		},
		{
			name:           "validation error",
			body:           `{"data":[{"name":""}]}`,
			vars:           map[string]string{"categoryID": "1"},
			wantCode:       http.StatusUnprocessableEntity,
			wantValidation: true,
		},
		{
			name:     "service batch empty",
			body:     `{"data":[{"name":"A"}]}`,
			vars:     map[string]string{"categoryID": "2"},
			svcErr:   attrDom.ErrBatchEmpty,
			wantCode: http.StatusBadRequest,
			wantMsg:  "no attributes provided for batch",
		},
		{
			name:     "service category not found",
			body:     `{"data":[{"name":"A"}]}`,
			vars:     map[string]string{"categoryID": "3"},
			svcErr:   catDom.ErrCategoryNotFound,
			wantCode: http.StatusNotFound,
			wantMsg:  "category not found",
		},
		{
			name:     "service conflict",
			body:     `{"data":[{"name":"A"}]}`,
			vars:     map[string]string{"categoryID": "4"},
			svcErr:   attrDom.ErrAttributeAlreadyExists,
			wantCode: http.StatusConflict,
			wantMsg:  "attribute already exists",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := mocks.NewMockAttributeService(t)
			h := newHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewReader([]byte(tc.body)))
			req = withChiParams(req, tc.vars)
			w := httptest.NewRecorder()

			if tc.svcErr != nil {
				mockSvc.
					On("CreateAttributesBatch", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.svcErr).
					Once()
			}

			h.CreateAttributesBatch(w, req)
			assert.Equal(t, tc.wantCode, w.Code)

			if tc.wantValidation {
				var resp httputils.ValidationErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Errors)
			} else {
				var resp httputils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tc.wantMsg, resp.Message)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestCreateAttribute_Success(t *testing.T) {
	mockSvc := mocks.NewMockAttributeService(t)
	h := newHandler(mockSvc)

	body := map[string]interface{}{"name": "AttrName", "unit": "u"}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/category/1/attribute", bytes.NewReader(raw))
	req = withChiParams(req, map[string]string{"categoryID": "1"})
	w := httptest.NewRecorder()

	mockSvc.
		On("CreateAttribute", mock.Anything, int64(1), mock.MatchedBy(func(a *attrDom.Attribute) bool {
			return a.Name == "AttrName" && a.Unit != nil && *a.Unit == "u"
		})).
		Return(&attrDom.Attribute{ID: 55, Name: "AttrName", Unit: strPtr("u"), CategoryID: 1}, nil).
		Once()

	h.CreateAttribute(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/api/admin/category/1/attribute/55", w.Header().Get("Location"))
	var resp dto.AttributeResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(55), resp.ID)
	assert.Equal(t, "AttrName", resp.Name)
	assert.Equal(t, int64(1), resp.CategoryID)
	mockSvc.AssertExpectations(t)
}

func TestCreateAttribute_ValidationError(t *testing.T) {
	mockSvc := mocks.NewMockAttributeService(t)
	h := newHandler(mockSvc)

	body := map[string]interface{}{"name": ""}
	raw, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/category/1/attribute", bytes.NewReader(raw))
	req = withChiParams(req, map[string]string{"categoryID": "1"})
	w := httptest.NewRecorder()

	h.CreateAttribute(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp httputils.ValidationErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Errors)
}

func strPtr(s string) *string { return &s }
