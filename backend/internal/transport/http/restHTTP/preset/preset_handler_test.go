package preset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"log/slog"

	domPreset "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-chi/chi/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// fakeValidator implements interfaces.Validator
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

func newHandler(mockSvc *mocks.MockPresetService, valErr error) *Handler {
	val := fakeValidator{err: valErr}
	deps, err := NewDeps(val, slog.New(slog.DiscardHandler), mockSvc)
	if err != nil {
		panic("failed to construct deps: " + err.Error())
	}
	return New(deps)
}

func withChiParam(r *http.Request, key, val string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func TestCreatePreset(t *testing.T) {
	validBody := map[string]interface{}{
		"name":        "MyPreset",
		"total_price": 42.5,
		"items": []map[string]interface{}{
			{"product_id": 1, "price": 42.5},
		},
	}
	bs, _ := json.Marshal(validBody)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", bytes.NewReader(bs))
		w := httptest.NewRecorder()

		mockSvc.
			On("Create", mock.Anything, mock.MatchedBy(func(p *domPreset.Preset) bool {
				return p.Name == "MyPreset" &&
					len(p.Items) == 1 &&
					p.TotalPrice == 42.5
			})).
			Return(int64(99), nil).
			Once()

		h.Create(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		assert.Equal(t, "/api/admin/presets/99", w.Header().Get("Location"))

		var resp dto.PresetResponseID
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, int64(99), resp.PresetID)
		assert.Equal(t, "Preset created successfully", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", bytes.NewReader([]byte("{bad json")))
		w := httptest.NewRecorder()

		h.Create(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var e httputils.ErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &e))
		assert.Contains(t, e.Message, "Invalid JSON")
	})

	t.Run("validation error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		errs := validator.ValidationErrors{fakeFieldError{}}
		h := newHandler(mockSvc, errs)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", bytes.NewReader(bs))
		w := httptest.NewRecorder()

		h.Create(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var verrs dto.ValidationErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &verrs))
		assert.Equal(t, "Validation failed", verrs.Message)
		assert.NotEmpty(t, verrs.Errors)
	})

	for name, serviceErr := range map[string]error{
		"name exists":        domPreset.ErrPresetAlreadyExists,
		"no items":           domPreset.ErrNoItems,
		"name too long":      domPreset.ErrNameTooLong,
		"price mismatch":     domPreset.ErrTotalPriceMismatch,
		"invalid product id": domPreset.ErrInvalidProductID,
	} {
		t.Run("service "+name, func(t *testing.T) {
			mockSvc := new(mocks.MockPresetService)
			h := newHandler(mockSvc, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", bytes.NewReader(bs))
			w := httptest.NewRecorder()

			mockSvc.
				On("Create", mock.Anything, mock.Anything).
				Return(int64(0), serviceErr).
				Once()

			h.Create(w, req)
			// all these map to 4xx
			assert.Contains(t, []int{http.StatusConflict, http.StatusUnprocessableEntity}, w.Code)
		})
	}

	t.Run("service internal error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/presets", bytes.NewReader(bs))
		w := httptest.NewRecorder()

		mockSvc.
			On("Create", mock.Anything, mock.Anything).
			Return(int64(0), errors.New("boom")).
			Once()

		h.Create(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetPreset(t *testing.T) {
	imageURL := "http://example.com/image.jpg"
	// prepare a sample domain.Preset
	d := &domPreset.Preset{
		ID:         5,
		Name:       "X",
		TotalPrice: 10,
		Items: []domPreset.PresetItem{
			{ProductID: 2, PresetID: 1, Product: &domProduct.ProductSummary{ID: 2, Name: "Product 2", Price: 2}},
			{ProductID: 3, PresetID: 1, Product: &domProduct.ProductSummary{ID: 3, Name: "Product 3", Price: 8, ImageURL: &imageURL}},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/presets/5", nil), "id", "5")
		w := httptest.NewRecorder()

		mockSvc.
			On("Get", mock.Anything, int64(5)).
			Return(d, nil).
			Once()

		h.Get(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp dto.PresetResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, d.ID, resp.PresetID)
		require.Len(t, resp.Items, 2)
		assert.Equal(t, int64(2), resp.Items[0].Product.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/presets/abc", nil), "id", "abc")
		w := httptest.NewRecorder()

		h.Get(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/presets/7", nil), "id", "7")
		w := httptest.NewRecorder()

		mockSvc.
			On("Get", mock.Anything, int64(7)).
			Return(nil, domPreset.ErrPresetNotFound).
			Once()

		h.Get(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/presets/8", nil), "id", "8")
		w := httptest.NewRecorder()

		mockSvc.
			On("Get", mock.Anything, int64(8)).
			Return(nil, errors.New("fail")).
			Once()

		h.Get(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeletePreset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/presets/3", nil), "id", "3")
		w := httptest.NewRecorder()

		mockSvc.
			On("Delete", mock.Anything, int64(3)).
			Return(nil).
			Once()

		h.Delete(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/presets/xyz", nil), "id", "xyz")
		w := httptest.NewRecorder()

		h.Delete(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/presets/42", nil), "id", "42")
		w := httptest.NewRecorder()

		mockSvc.
			On("Delete", mock.Anything, int64(42)).
			Return(domPreset.ErrPresetNotFound).
			Once()

		h.Delete(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/presets/99", nil), "id", "99")
		w := httptest.NewRecorder()

		mockSvc.
			On("Delete", mock.Anything, int64(99)).
			Return(errors.New("oops")).
			Once()

		h.Delete(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListDetailedAndShort(t *testing.T) {
	p1 := domPreset.Preset{ID: 1, Name: "A"}
	p2 := domPreset.Preset{ID: 2, Name: "B"}

	t.Run("ListDetailed success", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/presets/detailed", nil)
		w := httptest.NewRecorder()

		mockSvc.
			On("ListDetailed", mock.Anything).
			Return([]domPreset.Preset{p1, p2}, nil).
			Once()

		h.ListDetailed(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var out []dto.PresetResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
		assert.Len(t, out, 2)
	})

	t.Run("ListDetailed error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/presets/detailed", nil)
		w := httptest.NewRecorder()

		mockSvc.
			On("ListDetailed", mock.Anything).
			Return(nil, errors.New("fail")).
			Once()

		h.ListDetailed(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ListShort success", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/presets", nil)
		w := httptest.NewRecorder()

		mockSvc.
			On("ListShort", mock.Anything).
			Return([]domPreset.Preset{p1, p2}, nil).
			Once()

		h.ListShort(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var out []dto.PresetShortResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
		assert.Len(t, out, 2)
	})

	t.Run("ListShort error", func(t *testing.T) {
		mockSvc := new(mocks.MockPresetService)
		h := newHandler(mockSvc, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/presets", nil)
		w := httptest.NewRecorder()

		mockSvc.
			On("ListShort", mock.Anything).
			Return(nil, errors.New("fail")).
			Once()

		h.ListShort(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
