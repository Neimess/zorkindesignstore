package coefficients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/coefficients/dto"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockService struct{ mock.Mock }

func (m *MockService) Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockService) Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockService) List(ctx context.Context) ([]domCoeff.Coefficient, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domCoeff.Coefficient), args.Error(1)
}
func (m *MockService) Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type HandlerSuite struct {
	suite.Suite
	h  *Handler
	ms *MockService
}

func (s *HandlerSuite) SetupTest() {
	s.ms = new(MockService)
	deps, err := NewDeps(slog.Default(), s.ms)
	s.Require().NoError(err)
	s.h = New(deps)
}

func withChiParam(r *http.Request, key, val string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func (s *HandlerSuite) TestCreate() {
	body := dto.CoefficientRequest{Name: "A", Value: 1.1}
	b, _ := json.Marshal(body)
	s.ms.On("Create", mock.Anything, mock.Anything).Return(&domCoeff.Coefficient{ID: 1, Name: "A", Value: 1.1}, nil).Once()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/coefficients", bytes.NewReader(b))
	w := httptest.NewRecorder()
	s.h.Create(w, req)
	s.Equal(http.StatusCreated, w.Code)
	var resp dto.CoefficientResponse
	s.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	s.Equal("A", resp.Name)
}

func (s *HandlerSuite) TestGet() {
	s.ms.On("Get", mock.Anything, int64(1)).Return(&domCoeff.Coefficient{ID: 1, Name: "A", Value: 1.1}, nil).Once()
	req := withChiParam(httptest.NewRequest(http.MethodGet, "/api/admin/coefficients/1", nil), "id", "1")
	w := httptest.NewRecorder()
	s.h.Get(w, req)
	s.Equal(http.StatusOK, w.Code)
	var resp dto.CoefficientResponse
	s.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	s.Equal(int64(1), resp.ID)
}

func (s *HandlerSuite) TestList() {
	list := []domCoeff.Coefficient{{ID: 1, Name: "A", Value: 1.1}}
	s.ms.On("List", mock.Anything).Return(list, nil).Once()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/coefficients", nil)
	w := httptest.NewRecorder()
	s.h.List(w, req)
	s.Equal(http.StatusOK, w.Code)
	var resp []dto.CoefficientResponse
	s.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	s.Len(resp, 1)
}

func (s *HandlerSuite) TestUpdate() {
	body := dto.CoefficientRequest{Name: "B", Value: 2.2}
	b, _ := json.Marshal(body)
	s.ms.On("Update", mock.Anything, mock.Anything).Return(&domCoeff.Coefficient{ID: 2, Name: "B", Value: 2.2}, nil).Once()
	req := withChiParam(httptest.NewRequest(http.MethodPut, "/api/admin/coefficients/2", bytes.NewReader(b)), "id", "2")
	w := httptest.NewRecorder()
	s.h.Update(w, req)
	s.Equal(http.StatusOK, w.Code)
	var resp dto.CoefficientResponse
	s.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	s.Equal("B", resp.Name)
}

func (s *HandlerSuite) TestDelete() {
	s.ms.On("Delete", mock.Anything, int64(3)).Return(nil).Once()
	req := withChiParam(httptest.NewRequest(http.MethodDelete, "/api/admin/coefficients/3", nil), "id", "3")
	w := httptest.NewRecorder()
	s.h.Delete(w, req)
	s.Equal(http.StatusNoContent, w.Code)
}

func (s *HandlerSuite) TestServiceErrors() {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"not found", domCoeff.ErrCoefficientNotFound, http.StatusNotFound},
		{"conflict", domCoeff.ErrCoefficientAlreadyExists, http.StatusConflict},
		{"empty name", domCoeff.ErrEmptyName, http.StatusUnprocessableEntity},
		{"name too long", domCoeff.ErrNameTooLong, http.StatusUnprocessableEntity},
		{"internal", errors.New("fail"), http.StatusInternalServerError},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.ms.On("Create", mock.Anything, mock.Anything).Return((*domCoeff.Coefficient)(nil), tc.err).Once()
			body := dto.CoefficientRequest{Name: "X", Value: 1.1}
			b, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/admin/coefficients", bytes.NewReader(b))
			w := httptest.NewRecorder()
			s.h.Create(w, req)
			assert.Equal(s.T(), tc.want, w.Code)
		})
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
