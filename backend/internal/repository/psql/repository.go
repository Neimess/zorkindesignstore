package repository

type Repositories struct {
	ProductRepository  *ProductRepository
	CategoryRepository *CategoryRepository
}

func New(deps Deps) *Repositories {
	return &Repositories{
		ProductRepository: NewProductRepository(deps.DB, deps.Logger),
		CategoryRepository: NewCategoryRepository(deps.DB, deps.Logger),
	}
}
