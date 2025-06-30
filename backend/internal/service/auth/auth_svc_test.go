package auth_test

import (
	"errors"
	"fmt"
	"testing"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/service/auth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockJWTGen struct{ mock.Mock }

func (m *mockJWTGen) Generate(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

type AuthServiceSuite struct {
	suite.Suite
	svc  *auth.Service
	mock *mockJWTGen
	log  *slog.Logger
}

func (s *AuthServiceSuite) SetupTest() {
	s.mock = new(mockJWTGen)
	s.log = slog.Default()
	deps, _ := auth.NewDeps(s.mock, s.log)
	s.svc = auth.New(deps)
}

func (s *AuthServiceSuite) TestGenerateToken() {
	type testCase struct {
		name      string
		userID    string
		mockSetup func()
		expectErr bool
	}

	tests := []testCase{
		{
			name:   "success",
			userID: "user1",
			mockSetup: func() {
				s.mock.On("Generate", "user1").Return("token123", nil).Once()
			},
			expectErr: false,
		},
		{
			name:   "error from generator",
			userID: "user2",
			mockSetup: func() {
				s.mock.On("Generate", "user2").Return("", errors.New("fail")).Once()
			},
			expectErr: true,
		},
		{
			name:   "empty userID",
			userID: "",
			mockSetup: func() {
				s.mock.On("Generate", "").Return("", fmt.Errorf("empty userID")).Once()
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			token, err := s.svc.GenerateToken(tc.userID)
			if tc.expectErr {
				s.Error(err)
				s.Empty(token)
			} else {
				s.NoError(err)
				s.Equal("token123", token)
			}
			s.mock.AssertExpectations(s.T())
		})
	}
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
