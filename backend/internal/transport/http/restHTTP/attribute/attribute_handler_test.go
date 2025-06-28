package attribute

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"
)

type fakeValidator struct{ err error }

func (f fakeValidator) StructCtx(ctx context.Context, s interface{}) error { return f.err }

func newHandler(mockSvc *mocks.MockAttributeService, valErr error) *Handler {
	dep, err := NewDeps(mockSvc, validator.New())
	if err != nil {
		return nil
	}
	h := New(dep)
	h.val = fakeValidator{err: valErr}
	return h
}

func TestCreateAttributesBatch_Success(t *testing.T) {
	mockSvc := mocks.NewMockAttributeService(t)
	h := newHandler(mockSvc, nil)

	input := []map[string]interface{}{{"name": "A", "unit": "u"}}
	body, err := json.Marshal(map[string]interface{}{"data": input})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewReader(body))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("categoryID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
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
		name     string
		body     string
		vars     map[string]string
		valErr   error
		svcErr   error
		wantCode int
		wantMsg  string
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
			h := newHandler(mockSvc, nil)

			req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewReader([]byte(tc.body)))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("categoryID", tc.vars["categoryID"])
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			w := httptest.NewRecorder()

			if tc.svcErr != nil {
				mockSvc.
					On("CreateAttributesBatch", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.svcErr).
					Once()
			}

			h.CreateAttributesBatch(w, req)
			assert.Equal(t, tc.wantCode, w.Code)

			var resp httputils.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantMsg, resp.Message)
			mockSvc.AssertExpectations(t)
		})
	}
}
