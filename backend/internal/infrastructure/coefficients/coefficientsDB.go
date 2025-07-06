package coefficients

import (
	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
)

type coefficientDB struct {
	ID    int64   `db:"coefficient_id"`
	Name  string  `db:"name"`
	Value float64 `db:"value"`
}

func (c coefficientDB) toDomain() *domCoeff.Coefficient {
	return &domCoeff.Coefficient{
		ID:    c.ID,
		Name:  c.Name,
		Value: c.Value,
	}
}

func rawCoeffListToDomain(raws []coefficientDB) []domCoeff.Coefficient {
	coeffs := make([]domCoeff.Coefficient, len(raws))
	for i, r := range raws {
		coeffs[i] = *r.toDomain()
	}
	return coeffs
}
