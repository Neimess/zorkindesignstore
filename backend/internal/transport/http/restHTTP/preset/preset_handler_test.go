package preset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"log/slog"

	domPreset "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset/mocks"
	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// === Helpers ===
type fakeValidator struct{ err error }

func (f fakeValidator) StructCtx(ctx context.Context, s interface{}) error {
	return f.err
}

type fakeFieldError struct{}

func (f fakeFieldError) Tag() string                       { return "required" }
func (f fakeFieldError) ActualTag() string                 { return "required" }
func (f fakeFieldError) Namespace() string                 { return "Preset.Name" }
func (f fakeFieldError) StructNamespace() string           { return "Preset.Name" }
func (f fakeFieldError) Field() string                     { return "Name" }
func (f fakeFieldError) StructField() string               { return "Name" }
func (f fakeFieldError) Value() interface{}                { return "" }
func (f fakeFieldError) Param() string                     { return "" }
func (f fakeFieldError) Kind() reflect.Kind                { return reflect.String }
func (f fakeFieldError) Type() reflect.Type                { return reflect.TypeOf("") }
func (f fakeFieldError) Translate(ut ut.Translator) string { return f.Error() }
func (f fakeFieldError) Error() string                     { return "Name is required" }

type PresetHandlerSuite struct {
	suite.Suite
	h       *Handler
	mockSvc *mocks.MockPresetService
	valErr  error
}

func (s *PresetHandlerSuite) SetupTest() {
	s.mockSvc = new(mocks.MockPresetService)
	val := fakeValidator{err: s.valErr}
	deps, err := NewDeps(val, slog.New(slog.DiscardHandler), s.mockSvc)
	s.Require().NoError(err)
	s.h = New(deps)
}

func withChiParam(r *http.Request, key, val string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, val)
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
	bs, _ := json.Marshal(validBody)

	type svcResult struct {
		p   *domPreset.Preset
		err error
	}

	tests := []struct {
		name           string
		validatorError error
		svcReturn      svcResult
		wantStatus     int
		wantLocation   bool
	}{
		{
			name:           "success",
			validatorError: nil,
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
			name:           "invalid JSON",
			validatorError: nil,
			svcReturn:      svcResult{},
			wantStatus:     http.StatusBadRequest,
		},
		{
			name:           "validation error",
			validatorError: validator.ValidationErrors{fakeFieldError{}},
			svcReturn:      svcResult{},
			wantStatus:     http.StatusUnprocessableEntity,
		},
		{
			name:           "domain errors → 4xx",
			validatorError: nil,
			svcReturn:      svcResult{p: nil, err: domPreset.ErrPresetAlreadyExists},
			wantStatus:     http.StatusConflict,
		},
		{
			name:           "internal service error",
			validatorError: nil,
			svcReturn:      svcResult{p: nil, err: errors.New("boom")},
			wantStatus:     http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			s.valErr = tc.validatorError
			val := fakeValidator{err: tc.validatorError}
			deps, err := NewDeps(val, slog.New(slog.DiscardHandler), s.mockSvc)
			s.Require().NoError(err)
			s.h = New(deps)

			reqBody := bytes.NewReader(bs)
			if tc.name == "invalid JSON" {
				reqBody = bytes.NewReader([]byte("{bad json"))
			}
			req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", reqBody)
			w := httptest.NewRecorder()

			if tc.name == "success" || tc.name == "domain errors → 4xx" || tc.name == "internal service error" {
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
				assert.Equal(s.T(), "/api/admin/presets/"+id, loc)
			}

			s.mockSvc.AssertExpectations(s.T())
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
		"total_price": 99.99,
		"items": []map[string]interface{}{
			{"product_id": 1},
			{"product_id": 2},
		},
	}
	bs, _ := json.Marshal(validBody)

	type svcResult struct {
		p   *domPreset.Preset
		err error
	}

	tests := []struct {
		name       string
		idParam    string
		validator  error
		svcReturn  svcResult
		wantStatus int
	}{
		{
			name:      "success",
			idParam:   "10",
			validator: nil,
			svcReturn: svcResult{
				p: &domPreset.Preset{
					ID:         10,
					Name:       "UpdatedPreset",
					TotalPrice: 99.99,
					Items: []domPreset.PresetItem{
						{ProductID: 1, PresetID: 10},
						{ProductID: 2, PresetID: 10},
					},
				},
				err: nil,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			idParam:    "abc",
			validator:  nil,
			svcReturn:  svcResult{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid JSON",
			idParam:    "10",
			validator:  nil,
			svcReturn:  svcResult{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "validation error",
			idParam:    "10",
			validator:  validator.ValidationErrors{fakeFieldError{}},
			svcReturn:  svcResult{},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "not found",
			idParam:    "10",
			validator:  nil,
			svcReturn:  svcResult{p: nil, err: domPreset.ErrPresetNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "internal error",
			idParam:    "10",
			validator:  nil,
			svcReturn:  svcResult{p: nil, err: errors.New("boom")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			s.valErr = tc.validator
			val := fakeValidator{err: tc.validator}
			deps, err := NewDeps(val, slog.New(slog.DiscardHandler), s.mockSvc)
			s.Require().NoError(err)
			s.h = New(deps)

			var reqBody *bytes.Reader
			if tc.name == "invalid JSON" {
				reqBody = bytes.NewReader([]byte("{bad json"))
			} else {
				reqBody = bytes.NewReader(bs)
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
