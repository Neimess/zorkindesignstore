package preset

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type PresetRepository interface {
	Create(ctx context.Context, p *preset.Preset) (*preset.Preset, error)
	Get(ctx context.Context, id int64) (*preset.Preset, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, p *preset.Preset) (*preset.Preset, error)
	ListDetailed(ctx context.Context) ([]preset.Preset, error)
	ListShort(ctx context.Context) ([]preset.Preset, error)
}

type Service struct {
	repo PresetRepository
	log  *slog.Logger
}

type Deps struct {
	Repo PresetRepository
	log  *slog.Logger
}

func NewDeps(repo PresetRepository, log *slog.Logger) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("preset: missing PresetRepository")
	}
	if log == nil {
		return nil, errors.New("preset: missing logger")
	}
	return &Deps{Repo: repo, log: log.With("component", "service.preset")}, nil
}

func New(d *Deps) *Service {
	return &Service{
		repo: d.Repo,
		log:  d.log,
	}
}

func (s *Service) Create(ctx context.Context, p *preset.Preset) (*preset.Preset, error) {
	const op = "service.preset.Create"
	log := s.log.With("op", op)

	if err := p.Validate(); err != nil {
		return nil, err
	}

	p, err := s.repo.Create(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrConflict: preset.ErrPresetAlreadyExists,
			der.ErrNotFound: preset.ErrInvalidProductID,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("preset created", slog.Int64("preset_id", p.ID))
	return p, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*preset.Preset, error) {
	const op = "service.preset.Get"
	log := s.log.With("op", op)

	p, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, preset.ErrPresetNotFound
		}
		log.Error("repo.Get failed", slog.Int64("preset_id", id), slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset retrieved", slog.Int64("preset_id", p.ID))
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.preset.Delete"
	log := s.log.With("op", op)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound: nil,
		}
		return utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("preset deleted", slog.Int64("preset_id", id))
	return nil
}

func (s *Service) ListDetailed(ctx context.Context) ([]preset.Preset, error) {
	const op = "service.preset.ListDetailed"
	log := s.log.With("op", op)

	list, err := s.repo.ListDetailed(ctx)
	if err != nil {
		log.Error("repo.ListDetailed failed", slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("presets retrieved", slog.Int("count", len(list)))
	return list, nil
}

func (s *Service) ListShort(ctx context.Context) ([]preset.Preset, error) {
	const op = "service.preset.ListShort"
	log := s.log.With("op", op)

	list, err := s.repo.ListShort(ctx)
	if err != nil {
		log.Error("repo.ListShort failed", slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("presets (short) retrieved", slog.Int("count", len(list)))
	return list, nil
}

func (s *Service) Update(ctx context.Context, p *preset.Preset) (*preset.Preset, error) {
	const op = "service.preset.Update"
	log := s.log.With("op", op)

	res, err := s.repo.Update(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrConflict: preset.ErrPresetAlreadyExists,
			der.ErrNotFound: preset.ErrPresetNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("preset updated", slog.Int64("preset_id", res.ID))
	return res, nil
}
