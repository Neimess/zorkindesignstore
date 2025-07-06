package coefficients

import (
	"context"
	"errors"
	"log/slog"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type CoefficientRepository interface {
	Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error)
	Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error)
	List(ctx context.Context) ([]domCoeff.Coefficient, error)
	Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo CoefficientRepository
	log  *slog.Logger
}

type Deps struct {
	Repo CoefficientRepository
	Log  *slog.Logger
}

func NewDeps(repo CoefficientRepository, log *slog.Logger) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("coefficients: missing repository")
	}
	if log == nil {
		return nil, errors.New("coefficients: missing logger")
	}
	return &Deps{Repo: repo, Log: log.With("component", "service.coefficients")}, nil
}

func New(d *Deps) *Service {
	return &Service{
		repo: d.Repo,
		log:  d.Log,
	}
}

func (s *Service) Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	const op = "service.coefficients.Create"
	log := s.log.With("op", op)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	res, err := s.repo.Create(ctx, c)
	if err != nil {
		mapping := map[error]error{
			der.ErrConflict: domCoeff.ErrCoefficientAlreadyExists,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error) {
	const op = "service.coefficients.Get"
	log := s.log.With("op", op)
	res, err := s.repo.Get(ctx, id)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domCoeff.ErrCoefficientNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *Service) List(ctx context.Context) ([]domCoeff.Coefficient, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	const op = "service.coefficients.Update"
	log := s.log.With("op", op)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	res, err := s.repo.Update(ctx, c)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domCoeff.ErrCoefficientNotFound,
			der.ErrConflict: domCoeff.ErrCoefficientAlreadyExists,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.coefficients.Delete"
	log := s.log.With("op", op)
	err := s.repo.Delete(ctx, id)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domCoeff.ErrCoefficientNotFound,
		}
		return utils.ErrorHandler(log, op, err, mapping)
	}
	return nil
}
