package service

import (
	"context"
	"errors"
	"log/slog"

	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type ServiceRepository interface {
	Create(ctx context.Context, s *domService.Service) (*domService.Service, error)
	Get(ctx context.Context, id int64) (*domService.Service, error)
	List(ctx context.Context) ([]domService.Service, error)
	Update(ctx context.Context, s *domService.Service) (*domService.Service, error)
	Delete(ctx context.Context, id int64) error
	AddServicesToProduct(ctx context.Context, productID int64, serviceIDs []int64) error
	GetServicesByProduct(ctx context.Context, productID int64) ([]domService.Service, error)
}

type ServiceSvc struct {
	repo ServiceRepository
	log  *slog.Logger
}

type Deps struct {
	Repo ServiceRepository
	Log  *slog.Logger
}

func NewDeps(repo ServiceRepository, log *slog.Logger) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("service: missing repository")
	}
	if log == nil {
		return nil, errors.New("service: missing logger")
	}
	return &Deps{Repo: repo, Log: log.With("component", "service.service")}, nil
}

func New(d *Deps) *ServiceSvc {
	return &ServiceSvc{
		repo: d.Repo,
		log:  d.Log,
	}
}

func (s *ServiceSvc) Create(ctx context.Context, serv *domService.Service) (*domService.Service, error) {
	const op = "service.service.Create"
	log := s.log.With("op", op)
	if serv.Name == "" {
		return nil, domService.ErrEmptyName
	}
	res, err := s.repo.Create(ctx, serv)
	if err != nil {
		mapping := map[error]error{
			der.ErrConflict: domService.ErrServiceAlreadyExists,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *ServiceSvc) Get(ctx context.Context, id int64) (*domService.Service, error) {
	const op = "service.service.Get"
	log := s.log.With("op", op)
	res, err := s.repo.Get(ctx, id)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *ServiceSvc) List(ctx context.Context) ([]domService.Service, error) {
	return s.repo.List(ctx)
}

func (s *ServiceSvc) Update(ctx context.Context, serv *domService.Service) (*domService.Service, error) {
	const op = "service.service.Update"
	log := s.log.With("op", op)
	if serv.Name == "" {
		return nil, domService.ErrEmptyName
	}
	res, err := s.repo.Update(ctx, serv)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
			der.ErrConflict: domService.ErrServiceAlreadyExists,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	return res, nil
}

func (s *ServiceSvc) Delete(ctx context.Context, id int64) error {
	const op = "service.service.Delete"
	log := s.log.With("op", op)
	err := s.repo.Delete(ctx, id)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
		}
		return utils.ErrorHandler(log, op, err, mapping)
	}
	return nil
}

func (s *ServiceSvc) AddServicesToProduct(ctx context.Context, productID int64, serviceIDs []int64) error {
	return s.repo.AddServicesToProduct(ctx, productID, serviceIDs)
}

func (s *ServiceSvc) GetServicesByProduct(ctx context.Context, productID int64) ([]domService.Service, error) {
	return s.repo.GetServicesByProduct(ctx, productID)
}
