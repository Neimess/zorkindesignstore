package coefficients_test

import (
	"context"
	"testing"

	"log/slog"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	coeffsvc "github.com/Neimess/zorkin-store-project/internal/service/coefficients"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockRepo) Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockRepo) List(ctx context.Context) ([]domCoeff.Coefficient, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domCoeff.Coefficient), args.Error(1)
}
func (m *MockRepo) Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*domCoeff.Coefficient), args.Error(1)
}
func (m *MockRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type CoefficientServiceSuite struct {
	suite.Suite
	svc  *coeffsvc.Service
	mock *MockRepo
}

func (s *CoefficientServiceSuite) SetupTest() {
	s.mock = new(MockRepo)
	deps, _ := coeffsvc.NewDeps(s.mock, slog.Default())
	s.svc = coeffsvc.New(deps)
}

func (s *CoefficientServiceSuite) TestCreate() {
	c := &domCoeff.Coefficient{Name: "A", Value: 1.1}
	s.mock.On("Create", mock.Anything, c).Return(c, nil).Once()
	res, err := s.svc.Create(context.Background(), c)
	s.NoError(err)
	s.Equal("A", res.Name)
}

func (s *CoefficientServiceSuite) TestCreate_Validation() {
	c := &domCoeff.Coefficient{Name: "", Value: 1.1}
	res, err := s.svc.Create(context.Background(), c)
	s.Error(err)
	s.Nil(res)
}

func (s *CoefficientServiceSuite) TestGet() {
	c := &domCoeff.Coefficient{ID: 1, Name: "A", Value: 1.1}
	s.mock.On("Get", mock.Anything, int64(1)).Return(c, nil).Once()
	res, err := s.svc.Get(context.Background(), 1)
	s.NoError(err)
	s.Equal(int64(1), res.ID)
}

func (s *CoefficientServiceSuite) TestList() {
	list := []domCoeff.Coefficient{{ID: 1, Name: "A", Value: 1.1}}
	s.mock.On("List", mock.Anything).Return(list, nil).Once()
	res, err := s.svc.List(context.Background())
	s.NoError(err)
	s.Len(res, 1)
}

func (s *CoefficientServiceSuite) TestUpdate() {
	c := &domCoeff.Coefficient{ID: 1, Name: "A", Value: 1.1}
	s.mock.On("Update", mock.Anything, c).Return(c, nil).Once()
	res, err := s.svc.Update(context.Background(), c)
	s.NoError(err)
	s.Equal("A", res.Name)
}

func (s *CoefficientServiceSuite) TestDelete() {
	s.mock.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
	err := s.svc.Delete(context.Background(), 1)
	s.NoError(err)
}

func TestCoefficientServiceSuite(t *testing.T) {
	suite.Run(t, new(CoefficientServiceSuite))
}
