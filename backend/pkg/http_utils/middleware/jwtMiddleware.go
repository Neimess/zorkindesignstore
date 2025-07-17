package middleware

import (
	"context"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type JWTCfg struct {
	Secret    []byte
	Algorithm string
	Issuer    string
	Audience  string
}

func NewJWTMiddleware(cfg JWTCfg, opts ...jwtmiddleware.Option) (*jwtmiddleware.JWTMiddleware, error) {
	keyFunc := func(ctx context.Context) (interface{}, error) {
		return []byte(cfg.Secret), nil
	}
	v, err := validator.New(
		keyFunc,
		validator.SignatureAlgorithm(cfg.Algorithm),
		cfg.Issuer,
		[]string{cfg.Audience},
		validator.WithCustomClaims(nil),
	)
	if err != nil {
		return nil, err
	}

	return jwtmiddleware.New(v.ValidateToken, opts...), nil
}
