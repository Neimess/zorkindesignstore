// internal/transport/http/restHTTP/auth/auth_handler_test.go
package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/auth/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/auth/mocks"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		mockToken      string
		mockErr        error
		wantStatus     int
		wantBodyStruct interface{}
	}{
		{
			name:       "success",
			mockToken:  "abc123",
			mockErr:    nil,
			wantStatus: http.StatusCreated,
			wantBodyStruct: &dto.TokenResponse{
				Token: "abc123",
			},
		},
		{
			name:       "service error",
			mockToken:  "",
			mockErr:    errors.New("something went wrong"),
			wantStatus: http.StatusInternalServerError,
			wantBodyStruct: &httputils.ErrorResponse{
				Message: "failed to generate token",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mocks.MockAuthService)
			handler := New(Deps{
				srv: mockSvc,
				log: slog.Default(),
			})

			mockSvc.
				On("GenerateToken", "ADMIN").
				Return(tc.mockToken, tc.mockErr).
				Once()

			req := httptest.NewRequest(http.MethodGet, "/api/admin/auth/secret-key", nil)
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			switch want := tc.wantBodyStruct.(type) {
			case *dto.TokenResponse:
				var got dto.TokenResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, want, &got)

			case *httputils.ErrorResponse:
				var got httputils.ErrorResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				assert.Equal(t, want, &got)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
