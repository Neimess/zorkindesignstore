package service

import (
	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
)

type ServiceDB struct {
	ID          int64   `db:"service_id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Price       float64 `db:"price"`
}

func (s ServiceDB) toDomain() *domService.Service {
	return &domService.Service{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Price:       s.Price,
	}
}

func rawServiceListToDomain(raws []ServiceDB) []domService.Service {
	services := make([]domService.Service, len(raws))
	for i, r := range raws {
		services[i] = *r.toDomain()
	}
	return services
}
