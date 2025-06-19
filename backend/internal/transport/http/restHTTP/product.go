package restHTTP

type ProductService interface {
	// Create(ctx context.Context, product *domain.Product, images []domain.ProductImage) error
	// GetByID(ctx context.Context, id int64) (*domain.ProductWithImages, error)
	// GetDetailed(ctx context.Context, id int64) (*domain.ProductWithCharacteristics, error)
	// GetAll(ctx context.Context) ([]domain.ProductWithImages, error)
}

type ProductHandler struct{}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{}
}
