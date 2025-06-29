package auth

import (
	"fmt"
	"log/slog"
)

type JWTGenerator interface {
	Generate(userID string) (string, error)
}

type Service struct {
	log *slog.Logger
	gen JWTGenerator
}

type Deps struct {
	jwtGen JWTGenerator
	log    *slog.Logger
}

func NewDeps(jwtGen JWTGenerator, log *slog.Logger) (Deps, error) {
	if jwtGen == nil {
		return Deps{}, fmt.Errorf("auth service: JWTGenerator is nil")
	}
	if log == nil {
		return Deps{}, fmt.Errorf("auth service: logger is nil")
	}
	return Deps{jwtGen: jwtGen, log: log.With("component", "service.auth")}, nil
}

func New(deps Deps) *Service {
	return &Service{
		log: deps.log,
		gen: deps.jwtGen,
	}
}

func (s *Service) GenerateToken(userID string) (string, error) {
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
