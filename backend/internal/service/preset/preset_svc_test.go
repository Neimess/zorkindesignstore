package preset_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	presetservice "github.com/Neimess/zorkin-store-project/internal/service/preset"
	"github.com/Neimess/zorkin-store-project/internal/service/preset/mocks"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PresetServiceSuite struct {
	suite.Suite
	svc      *presetservice.Service
	mockRepo *mocks.MockPresetRepository
	logger   *slog.Logger
}

func (s *PresetServiceSuite) SetupTest() {
	s.mockRepo = new(mocks.MockPresetRepository)
	s.logger = slog.Default()
	deps, _ := presetservice.NewDeps(s.mockRepo, s.logger)
	s.svc = presetservice.New(deps)
}

func validPreset() *preset.Preset {
	desc := "desc"
	img := "img"
	return &preset.Preset{
		ID:          1,
		Name:        "Test",
		Description: &desc,
		TotalPrice:  100.0,
		ImageURL:    &img,
		CreatedAt:   time.Now(),
		Items:       []preset.PresetItem{{ID: 1, ProductID: 1}},
	}
}

func (s *PresetServiceSuite) TestCreate() {
	type testCase struct {
		name      string
		input     *preset.Preset
		mockSetup func()
		expectErr error
	}

	tests := []testCase{
		{
			name:  "success",
			input: validPreset(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(validPreset(), nil).Once()
			},
			expectErr: nil,
		},
		{
			name:      "validation error (empty name)",
			input:     func() *preset.Preset { p := validPreset(); p.Name = ""; return p }(),
			mockSetup: func() {},
			expectErr: preset.ErrEmptyName,
		},
		{
			name:  "repo conflict error",
			input: validPreset(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, der.ErrConflict).Once()
			},
			expectErr: preset.ErrPresetAlreadyExists,
		},
		{
			name:  "repo not found error",
			input: validPreset(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: preset.ErrInvalidProductID,
		},
		{
			name:  "repo unknown error",
			input: validPreset(),
			mockSetup: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			res, err := s.svc.Create(context.Background(), tc.input)
			if tc.expectErr == nil {
				s.NoError(err)
				s.NotNil(res)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *PresetServiceSuite) TestGet() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expectErr error
	}

	p := validPreset()

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("Get", mock.Anything, int64(1)).Return(p, nil).Once()
			},
			expectErr: nil,
		},
		{
			name: "not found",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("Get", mock.Anything, int64(2)).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: preset.ErrPresetNotFound,
		},
		{
			name: "repo error",
			id:   3,
			mockSetup: func() {
				s.mockRepo.On("Get", mock.Anything, int64(3)).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			res, err := s.svc.Get(context.Background(), tc.id)
			if tc.expectErr == nil {
				s.NoError(err)
				s.NotNil(res)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *PresetServiceSuite) TestDelete() {
	type testCase struct {
		name      string
		id        int64
		mockSetup func()
		expectErr error
	}

	tests := []testCase{
		{
			name: "success",
			id:   1,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectErr: nil,
		},
		{
			name: "repo error",
			id:   2,
			mockSetup: func() {
				s.mockRepo.On("Delete", mock.Anything, int64(2)).Return(errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			err := s.svc.Delete(context.Background(), tc.id)
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *PresetServiceSuite) TestUpdate() {
	type testCase struct {
		name      string
		input     *preset.Preset
		mockSetup func()
		expectErr error
	}

	p := validPreset()

	tests := []testCase{
		{
			name:  "success",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(p, nil).Once()
			},
			expectErr: nil,
		},
		{
			name:  "repo conflict error",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, der.ErrConflict).Once()
			},
			expectErr: preset.ErrPresetAlreadyExists,
		},
		{
			name:  "repo not found error",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, der.ErrNotFound).Once()
			},
			expectErr: preset.ErrPresetNotFound,
		},
		{
			name:  "repo unknown error",
			input: p,
			mockSetup: func() {
				s.mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*preset.Preset")).Return(nil, errors.New("db fail")).Once()
			},
			expectErr: errors.New("db fail"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			res, err := s.svc.Update(context.Background(), tc.input)
			if tc.expectErr == nil {
				s.NoError(err)
				s.NotNil(res)
			} else {
				s.Error(err)
			}
		})
	}
}

func (s *PresetServiceSuite) TestListDetailed() {
	type testCase struct {
		name      string
		mockSetup func()
		expect    []preset.Preset
		expectErr bool
	}

	p := validPreset()

	tests := []testCase{
		{
			name: "success",
			mockSetup: func() {
				s.mockRepo.On("ListDetailed", mock.Anything).Return([]preset.Preset{*p}, nil).Once()
			},
			expect:    []preset.Preset{*p},
			expectErr: false,
		},
		{
			name: "empty list",
			mockSetup: func() {
				s.mockRepo.On("ListDetailed", mock.Anything).Return([]preset.Preset{}, nil).Once()
			},
			expect:    []preset.Preset{},
			expectErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				s.mockRepo.On("ListDetailed", mock.Anything).Return(nil, errors.New("db fail")).Once()
			},
			expect:    nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			got, err := s.svc.ListDetailed(context.Background())
			if tc.expectErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tc.expect, got)
			}
		})
	}
}

func (s *PresetServiceSuite) TestListShort() {
	type testCase struct {
		name      string
		mockSetup func()
		expect    []preset.Preset
		expectErr bool
	}

	p := validPreset()

	tests := []testCase{
		{
			name: "success",
			mockSetup: func() {
				s.mockRepo.On("ListShort", mock.Anything).Return([]preset.Preset{*p}, nil).Once()
			},
			expect:    []preset.Preset{*p},
			expectErr: false,
		},
		{
			name: "empty list",
			mockSetup: func() {
				s.mockRepo.On("ListShort", mock.Anything).Return([]preset.Preset{}, nil).Once()
			},
			expect:    []preset.Preset{},
			expectErr: false,
		},
		{
			name: "repo error",
			mockSetup: func() {
				s.mockRepo.On("ListShort", mock.Anything).Return(nil, errors.New("db fail")).Once()
			},
			expect:    nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()
			got, err := s.svc.ListShort(context.Background())
			if tc.expectErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tc.expect, got)
			}
		})
	}
}

func TestPresetServiceSuite(t *testing.T) {
	suite.Run(t, new(PresetServiceSuite))
}
