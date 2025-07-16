package preset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"log/slog"

	domPreset "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PresetHandlerSuite struct {
	suite.Suite
	h       *Handler
	mockSvc *mocks.MockPresetService
}

func (s *PresetHandlerSuite) SetupTest() {
	s.mockSvc = new(mocks.MockPresetService)
	deps, err := NewDeps(slog.New(slog.DiscardHandler), s.mockSvc)
	s.Require().NoError(err)
	s.h = New(deps)
}

func withChiParam(r *http.Request, key, val string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for k, v := range params {
		chiCtx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

// === TestCreatePreset ===

func (s *PresetHandlerSuite) TestCreatePreset() {
	validBody := map[string]interface{}{
		"name":        "MyPreset",
		"total_price": 42.5,
		"items": []map[string]interface{}{
			{"product_id": 1, "price": 42.5},
		},
	}

	type svcResult struct {
		p   *domPreset.Preset
		err error
	}

	tests := []struct {
		name         string
		body         interface{}
		svcReturn    svcResult
		wantStatus   int
		wantLocation bool
	}{
		{
			name: "success",
			body: validBody,
			svcReturn: svcResult{
				p: &domPreset.Preset{
					ID:         42,
					Name:       "MyPreset",
					TotalPrice: 42.5,
					Items: []domPreset.PresetItem{{
						ProductID: 1, PresetID: 42,
						Product: &domProduct.ProductSummary{ID: 1, Name: "Product 1", Price: 42.5},
					}},
				},
				err: nil,
			},
			wantStatus:   http.StatusCreated,
			wantLocation: true,
		},
		{
			name:       "invalid JSON",
			body:       `{`,
			svcReturn:  svcResult{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty name",
			body: map[string]interface{}{
				"name":        "",
				"total_price": 42.5,
				"items":       []map[string]interface{}{},
			},
			svcReturn:  svcResult{},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "domain errors → 4xx",
			body:       validBody,
			svcReturn:  svcResult{p: nil, err: domPreset.ErrPresetAlreadyExists},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "internal service error",
			body:       validBody,
			svcReturn:  svcResult{p: nil, err: errors.New("boom")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			var reqBody *bytes.Reader
			if s, ok := tc.body.(string); ok {
				reqBody = bytes.NewReader([]byte(s))
			} else {
				b, _ := json.Marshal(tc.body)
				reqBody = bytes.NewReader(b)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", reqBody)
			w := httptest.NewRecorder()

			shouldMock := tc.wantStatus == http.StatusCreated || tc.wantStatus == http.StatusConflict || tc.wantStatus == http.StatusInternalServerError
			if shouldMock {
				s.mockSvc.
					On("Create", mock.Anything, mock.Anything).
					Return(tc.svcReturn.p, tc.svcReturn.err).
					Once()
			}

			s.h.Create(w, req)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.wantLocation {
				loc := w.Header().Get("Location")
				id := strconv.Itoa(int(tc.svcReturn.p.ID))
				assert.Equal(s.T(), "/api/presets/"+id, loc)
			}

			if shouldMock {
				s.mockSvc.AssertExpectations(s.T())
			}
		})
	}
}

// === TestGetPreset ===

func (s *PresetHandlerSuite) TestGetPreset() {
	sample := &domPreset.Preset{
		ID:         5,
		Name:       "X",
		TotalPrice: 10,
		Items: []domPreset.PresetItem{
			{ProductID: 2, PresetID: 5, Product: &domProduct.ProductSummary{ID: 2, Name: "P2", Price: 2}},
			{ProductID: 3, PresetID: 5, Product: &domProduct.ProductSummary{ID: 3, Name: "P3", Price: 8}},
		},
	}

	tests := []struct {
		name       string
		param      string
		svcReturn  *domPreset.Preset
		svcErr     error
		wantStatus int
	}{
		{"success", "5", sample, nil, http.StatusOK},
		{"invalid id", "abc", nil, nil, http.StatusBadRequest},
		{"not found", "7", nil, domPreset.ErrPresetNotFound, http.StatusNotFound},
		{"service error", "8", nil, errors.New("fail"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/presets/"+tc.param, nil), "id", tc.param)
			w := httptest.NewRecorder()

			if tc.name == "success" || tc.name == "not found" || tc.name == "service error" {
				s.mockSvc.
					On("Get", mock.Anything, mock.MatchedBy(func(id int64) bool {
						i, _ := strconv.ParseInt(tc.param, 10, 64)
						return id == i
					})).
					Return(tc.svcReturn, tc.svcErr).
					Once()
			}

			s.h.Get(w, req)
			assert.Equal(s.T(), tc.wantStatus, w.Code)
			s.mockSvc.AssertExpectations(s.T())
		})
	}
}

// === TestDeletePreset ===

func (s *PresetHandlerSuite) TestDeletePreset() {
	tests := []struct {
		name            string
		param           string
		svcErr          error
		wantStatus      int
		checkIdempotent bool
	}{
		{"success", "3", nil, http.StatusNoContent, false},
		{"invalid id", "xyz", nil, http.StatusBadRequest, false},
		{"not found", "42", domPreset.ErrPresetNotFound, http.StatusNotFound, false},
		{"service error", "99", errors.New("oops"), http.StatusInternalServerError, false},
		{"idempotent retry", "100", nil, http.StatusNoContent, true},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/presets/"+tc.param, nil), "id", tc.param)
			w := httptest.NewRecorder()

			if tc.wantStatus != http.StatusBadRequest {
				for i := 0; i < 1+boolToInt(tc.checkIdempotent); i++ {
					s.mockSvc.
						On("Delete", mock.Anything, mock.MatchedBy(func(id int64) bool {
							i, _ := strconv.ParseInt(tc.param, 10, 64)
							return id == i
						})).
						Return(tc.svcErr).
						Once()
				}
			}

			s.h.Delete(w, req)
			assert.Equal(s.T(), tc.wantStatus, w.Code)

			if tc.checkIdempotent {
				w2 := httptest.NewRecorder()
				s.h.Delete(w2, req)
				assert.Equal(s.T(), tc.wantStatus, w2.Code)
			}

			s.mockSvc.AssertExpectations(s.T())
		})
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// === TestListDetailed and ListShort ===

func (s *PresetHandlerSuite) TestListPresetEndpoints() {
	p1 := domPreset.Preset{ID: 1, Name: "A"}
	p2 := domPreset.Preset{ID: 2, Name: "B"}

	tests := []struct {
		name        string
		method      string
		target      string
		mockMethod  string
		svcReturn   interface{}
		svcErr      error
		wantStatus  int
		wantPayload int
	}{
		{"ListDetailed OK", http.MethodGet, "/api/presets/detailed", "ListDetailed", []domPreset.Preset{p1, p2}, nil, http.StatusOK, 2},
		{"ListDetailed Err", http.MethodGet, "/api/presets/detailed", "ListDetailed", nil, errors.New("fail"), http.StatusInternalServerError, 0},
		{"ListShort OK", http.MethodGet, "/api/presets", "ListShort", []domPreset.Preset{p1, p2}, nil, http.StatusOK, 2},
		{"ListShort Err", http.MethodGet, "/api/presets", "ListShort", nil, errors.New("fail"), http.StatusInternalServerError, 0},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			req := httptest.NewRequest(tc.method, tc.target, nil)
			w := httptest.NewRecorder()

			if tc.wantStatus == http.StatusOK || tc.wantStatus == http.StatusInternalServerError {
				s.mockSvc.
					On(tc.mockMethod, mock.Anything).
					Return(tc.svcReturn, tc.svcErr).
					Once()
			}

			switch tc.mockMethod {
			case "ListDetailed":
				s.h.ListDetailed(w, req)
			case "ListShort":
				s.h.ListShort(w, req)
			}

			assert.Equal(s.T(), tc.wantStatus, w.Code)
			if tc.wantStatus == http.StatusOK {
				var arr []interface{}
				require.NoError(s.T(), json.Unmarshal(w.Body.Bytes(), &arr))
				assert.Len(s.T(), arr, tc.wantPayload)
			}

			s.mockSvc.AssertExpectations(s.T())
		})
	}
}

func (s *PresetHandlerSuite) TestUpdatePreset() {
	validBody := map[string]interface{}{
		"name":        "UpdatedPreset",
		"total_price": 50.0,
		"items": []map[string]interface{}{
			{"product_id": 2, "price": 50.0},
		},
	}

	type svcResult struct {
		p   *domPreset.Preset
		err error
	}

	tests := []struct {
		name       string
		idParam    string
		body       interface{}
		svcReturn  svcResult
		wantStatus int
	}{
		{
			name:    "success",
			idParam: "10",
			body:    validBody,
			svcReturn: svcResult{
				p: &domPreset.Preset{
					ID:         10,
					Name:       "UpdatedPreset",
					TotalPrice: 50.0,
					Items: []domPreset.PresetItem{{
						ProductID: 2, PresetID: 10,
						Product: &domProduct.ProductSummary{ID: 2, Name: "Product 2", Price: 50.0},
					}},
				},
				err: nil,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid JSON",
			idParam:    "10",
			body:       `{`,
			svcReturn:  svcResult{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "validation error - empty name",
			idParam: "10",
			body: map[string]interface{}{
				"name":        "",
				"total_price": 50.0,
				"items":       []map[string]interface{}{},
			},
			svcReturn:  svcResult{},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "not found",
			idParam:    "10",
			body:       validBody,
			svcReturn:  svcResult{p: nil, err: domPreset.ErrPresetNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "internal error",
			idParam:    "10",
			body:       validBody,
			svcReturn:  svcResult{p: nil, err: errors.New("boom")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			var reqBody *bytes.Reader
			if s, ok := tc.body.(string); ok {
				reqBody = bytes.NewReader([]byte(s))
			} else {
				b, _ := json.Marshal(tc.body)
				reqBody = bytes.NewReader(b)
			}

			req := withChiParam(
				httptest.NewRequest(http.MethodPut, "/api/admin/presets/"+tc.idParam, reqBody),
				"id", tc.idParam,
			)
			w := httptest.NewRecorder()

			if tc.name == "success" || tc.name == "not found" || tc.name == "internal error" {
				s.mockSvc.
					On("Update", mock.Anything, mock.Anything).
					Return(tc.svcReturn.p, tc.svcReturn.err).
					Once()
			}

			s.h.Update(w, req)

			assert.Equal(s.T(), tc.wantStatus, w.Code)
			s.mockSvc.AssertExpectations(s.T())
		})
	}
}

func TestPresetHandlerSuite(t *testing.T) {
	suite.Run(t, new(PresetHandlerSuite))
}

func newHandler(mockSvc *mocks.MockPresetService) *Handler {
	dep, err := NewDeps(slog.New(slog.DiscardHandler), mockSvc)
	if err != nil {
		return nil
	}
	h := New(dep)
	return h
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockPresetService(t)
			shouldMock := tc.svcErr != nil
			if shouldMock {
				mockSvc.
					EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(&domPreset.Preset{}, tc.svcErr).
					Once()
			}

			h := newHandler(mockSvc)
			req := httptest.NewRequest(http.MethodPost, "/presets", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			h.Create(w, req)

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
			if shouldMock {
				mockSvc.AssertExpectations(t)
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
			name:     "invalid preset id",
			body:     `{"name":"test","description":"test"}`,
			vars:     map[string]string{"id": "abc"},
			wantCode: http.StatusBadRequest,
			wantMsg:  "invalid id",
		},
		{
			name:           "validation error - empty name",
			body:           `{"name":"","description":"test"}`,
			vars:           map[string]string{"id": "1"},
			wantCode:       http.StatusUnprocessableEntity,
			wantValidation: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := mocks.NewMockPresetService(t)
			if tc.svcErr != nil {
				mockSvc.
					EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(&domPreset.Preset{}, tc.svcErr).
					Once()
			}

			h := newHandler(mockSvc)
			req := httptest.NewRequest(http.MethodPut, "/presets/1", strings.NewReader(tc.body))
			req = withChiParams(req, tc.vars)
			w := httptest.NewRecorder()

			h.Update(w, req)

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
