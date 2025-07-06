package dto

import domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"

func MapToDomain(r *CoefficientRequest) *domCoeff.Coefficient {
	return &domCoeff.Coefficient{
		Name:  r.Name,
		Value: r.Value,
	}
}

func MapToResponse(c *domCoeff.Coefficient) *CoefficientResponse {
	return &CoefficientResponse{
		ID:    c.ID,
		Name:  c.Name,
		Value: c.Value,
	}
}

func MapToResponseList(list []domCoeff.Coefficient) []CoefficientResponse {
	resp := make([]CoefficientResponse, len(list))
	for i, c := range list {
		resp[i] = *MapToResponse(&c)
	}
	return resp
}
