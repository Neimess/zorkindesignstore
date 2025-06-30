// internal/service/attribute/service_test.go
package attribute

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	mocks "github.com/Neimess/zorkin-store-project/internal/service/attribute/mocks"
	catMocks "github.com/Neimess/zorkin-store-project/internal/service/category/mocks"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	repoAttr *mocks.MockAttributeRepository
	repoCat  *catMocks.MockCategoryRepository
	svc      *Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.repoAttr = new(mocks.MockAttributeRepository)
	s.repoCat = new(catMocks.MockCategoryRepository)

	deps, err := NewDeps(s.repoAttr, s.repoCat, slog.New(slog.DiscardHandler))
	s.Require().NoError(err)
	s.svc = New(deps)
}

func (s *ServiceTestSuite) stubCatFound(id int64) {
	s.repoCat.
		On("GetByID", mock.Anything, id).
		Return(&catDom.Category{ID: id, Name: fmt.Sprintf("cat-%d", id)}, nil).
		Once()
}

func (s *ServiceTestSuite) stubCatNotFound(id int64) {
	s.repoCat.
		On("GetByID", mock.Anything, id).
		Return(nil, der.ErrNotFound).
		Once()
}

// ─── CreateAttributesBatch ────────────────────────────────────────────────────

func (s *ServiceTestSuite) TestCreateAttributesBatch() {
	batch := []attrDom.Attribute{
		{Name: "A1", CategoryID: 1},
		{Name: "A2", CategoryID: 1},
	}

	cases := []struct {
		name    string
		catID   int64
		attrs   []attrDom.Attribute
		setup   func()
		wantErr error
	}{
		{
			name:  "success",
			catID: 1,
			attrs: batch,
			setup: func() {
				s.stubCatFound(1)
				s.repoAttr.
					On("SaveBatch", mock.Anything, batch).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:    "empty slice",
			catID:   1,
			attrs:   []attrDom.Attribute{},
			setup:   func() {},
			wantErr: attrDom.ErrBatchEmpty,
		},
		{
			name:  "category not found",
			catID: 999,
			attrs: batch,
			setup: func() {
				s.stubCatNotFound(999)
			},
			wantErr: catDom.ErrCategoryNotFound,
		},
		{
			name:  "saving conflict",
			catID: 1,
			attrs: batch,
			setup: func() {
				s.stubCatFound(1)
				s.repoAttr.
					On("SaveBatch", mock.Anything, batch).
					Return(fmt.Errorf("%w", der.ErrConflict)).
					Once()
			},
			wantErr: attrDom.ErrAttributeAlreadyExists,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.T().Parallel()
			tc.setup()
			err := s.svc.CreateAttributesBatch(context.Background(), tc.catID, tc.attrs)

			if tc.wantErr != nil {
				assert.ErrorIs(s.T(), err, tc.wantErr)
			} else {
				assert.NoError(s.T(), err)
			}

			s.repoCat.AssertExpectations(s.T())
			s.repoAttr.AssertExpectations(s.T())
		})
	}
}

// ─── UpdateAttribute ─────────────────────────────────────────────────────---

func (s *ServiceTestSuite) TestUpdateAttribute() {
	goodAttr := &attrDom.Attribute{ID: 1, Name: "Attr1", CategoryID: 2}
	updatedAttr := &attrDom.Attribute{ID: 1, Name: "Attr1-upd", CategoryID: 2}
	badAttr := &attrDom.Attribute{ID: 2, Name: "", CategoryID: 2}

	cases := []struct {
		name    string
		input   *attrDom.Attribute
		setup   func()
		want    *attrDom.Attribute
		wantErr bool
	}{
		{
			name:  "success",
			input: &attrDom.Attribute{ID: 1, Name: "Attr1-upd", CategoryID: 2},
			setup: func() {
				s.stubCatFound(2)
				s.repoAttr.On("Update", mock.Anything, mock.MatchedBy(func(a *attrDom.Attribute) bool {
					return a.ID == 1 && a.Name == "Attr1-upd" && a.CategoryID == 2
				})).Return(updatedAttr, nil).Once()
			},
			want:    updatedAttr,
			wantErr: false,
		},
		{
			name:    "category not found",
			input:   goodAttr,
			setup:   func() { s.stubCatNotFound(2) },
			want:    nil,
			wantErr: true,
		},
		{
			name:  "repo error",
			input: goodAttr,
			setup: func() {
				s.stubCatFound(2)
				s.repoAttr.On("Update", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db fail")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "validation error",
			input:   badAttr,
			setup:   func() {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			tc.setup()
			res, err := s.svc.UpdateAttribute(context.Background(), tc.input)
			if tc.wantErr {
				assert.Error(s.T(), err)
				s.Nil(res)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.want, res)
			}
			s.repoCat.AssertExpectations(s.T())
			s.repoAttr.AssertExpectations(s.T())
		})
	}
}

// ─── DeleteAttribute ────────────────────────────────────────────────────────

func (s *ServiceTestSuite) TestDeleteAttribute() {
	cases := []struct {
		name    string
		id      int64
		setup   func()
		wantErr bool
	}{
		{
			name:    "success",
			id:      1,
			setup:   func() { s.repoAttr.On("Delete", mock.Anything, int64(1)).Return(nil).Once() },
			wantErr: false,
		},
		{
			name:    "not found (idempotent)",
			id:      2,
			setup:   func() { s.repoAttr.On("Delete", mock.Anything, int64(2)).Return(der.ErrNotFound).Once() },
			wantErr: false,
		},
		{
			name:    "repo error",
			id:      3,
			setup:   func() { s.repoAttr.On("Delete", mock.Anything, int64(3)).Return(fmt.Errorf("db fail")).Once() },
			wantErr: true,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			tc.setup()
			err := s.svc.DeleteAttribute(context.Background(), tc.id)
			if tc.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
			s.repoAttr.AssertExpectations(s.T())
		})
	}
}
