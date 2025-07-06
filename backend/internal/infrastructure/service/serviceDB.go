package service

import (
	"database/sql"

	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
)

type serviceDB struct {
	ID          int64          `db:"service_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Price       float64        `db:"price"`
}

func (s serviceDB) toDomain() *domService.Service {
	var desc *string
	if s.Description.Valid {
		desc = &s.Description.String
	}
	return &domService.Service{
		ID:          s.ID,
		Name:        s.Name,
		Description: desc,
		Price:       s.Price,
	}
}

func rawServiceListToDomain(raws []serviceDB) []domService.Service {
	services := make([]domService.Service, len(raws))
	for i, r := range raws {
		services[i] = *r.toDomain()
	}
	return services
}
