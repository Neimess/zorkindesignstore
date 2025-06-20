package service

import (
	"context"
	"testing"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repo "github.com/Neimess/zorkin-store-project/internal/repository/psql"
	"github.com/Neimess/zorkin-store-project/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestProductService_Create_Success(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	svc := NewProductService(mockRepo, silentLogger())

	ctx := context.Background()
	product := &domain.Product{
		Name:       "Бетонная кладка",
		Price:      570.1,
		CategoryID: 1,
		CreatedAt:  time.Now(),
	}

	mockRepo.EXPECT().
		Create(ctx, product).
		Return(int64(42), nil)

	id, err := svc.Create(ctx, product)
	require.NoError(t, err)
	require.Equal(t, int64(42), id)
}

func TestProductService_Create_NotFound(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	svc := NewProductService(mockRepo, silentLogger())
	ctx := context.Background()

	p := &domain.Product{Name: "Тест", Price: 123, CategoryID: 99}
	mockRepo.EXPECT().
		Create(ctx, p).
		Return(int64(0), repo.ErrProductNotFound)

	id, err := svc.Create(ctx, p)
	require.ErrorIs(t, err, repo.ErrProductNotFound)
	require.Zero(t, id)
}

func TestProductService_Create_UnicodeNames(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	svc := NewProductService(mockRepo, silentLogger())
	ctx := context.Background()

	products := []domain.Product{
		{Name: "Бетон", Price: 100, CategoryID: 1},
		{Name: "商品", Price: 100, CategoryID: 1},
		{Name: "📦🧱", Price: 100, CategoryID: 1},
		{Name: "المنتج", Price: 100, CategoryID: 1},
		{Name: "Товар 📱 产品", Price: 100, CategoryID: 1},
	}

	for _, p := range products {
		p := p
		mockRepo.EXPECT().
			Create(ctx, &p).
			Return(int64(1), nil)

		t.Run(p.Name, func(t *testing.T) {
			id, err := svc.Create(ctx, &p)
			require.NoError(t, err)
			require.Equal(t, int64(1), id)
		})
	}
}
