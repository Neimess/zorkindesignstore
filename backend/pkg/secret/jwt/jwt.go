package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	cfg JWTConfig
}

func NewGenerator(cfg JWTConfig) *JWTGenerator {
	return &JWTGenerator{cfg: cfg}
}

func (g *JWTGenerator) Generate(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    g.cfg.Issuer,
		Subject:   userID,
		Audience:  []string{g.cfg.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 1, 0)),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(g.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

type JWTConfig struct {
	Secret    string
	Algorithm string
	Issuer    string
	Audience  string
}
