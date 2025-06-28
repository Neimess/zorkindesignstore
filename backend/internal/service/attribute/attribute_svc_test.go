// internal/service/attribute/service_test.go
package attribute

import (
	"context"
	"fmt"
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

	deps, err := NewDeps(s.repoAttr, s.repoCat)
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
