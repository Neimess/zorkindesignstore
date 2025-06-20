package service

import (
	"fmt"
	"log/slog"

	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

type AuthService struct {
	log *slog.Logger
	gen JWTGenerator
}

func NewAuthService(jwtGen JWTGenerator, log *slog.Logger) *AuthService {
	if jwtGen == nil {
		panic("auth service: JWTGenerator is nil")
	}
	return &AuthService{
		log: logger.WithComponent(log, "service.auth"),
		gen: jwtGen,
	}
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	const op = "service.auth.GenerateToken"
	log := s.log.With("op", op)

	token, err := s.gen.Generate(userID)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token generated successfully", slog.String("user_id", userID))
	return token, nil
}
