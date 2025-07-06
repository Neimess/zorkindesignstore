package dto

import domService "github.com/Neimess/zorkin-store-project/internal/domain/service"

func MapToDomain(r *ServiceRequest) *domService.Service {
	return &domService.Service{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
	}
}

func MapToResponse(s *domService.Service) *ServiceResponse {
	return &ServiceResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Price:       s.Price,
	}
}

func MapToResponseList(list []domService.Service) []ServiceResponse {
	resp := make([]ServiceResponse, len(list))
	for i, s := range list {
		resp[i] = *MapToResponse(&s)
	}
	return resp
}
