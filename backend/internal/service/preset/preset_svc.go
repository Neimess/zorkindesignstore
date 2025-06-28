package preset

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type PresetRepository interface {
	Create(ctx context.Context, p *preset.Preset) (int64, error)
	Get(ctx context.Context, id int64) (*preset.Preset, error)
	Delete(ctx context.Context, id int64) error
	ListDetailed(ctx context.Context) ([]preset.Preset, error)
	ListShort(ctx context.Context) ([]preset.Preset, error)
}

type Service struct {
	repo PresetRepository
	log  *slog.Logger
}

type Deps struct {
	Repo PresetRepository
}

func NewDeps(repo PresetRepository) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("preset: missing PresetRepository")
	}
	return &Deps{Repo: repo}, nil
}

func New(d *Deps) *Service {
	return &Service{
		repo: d.Repo,
		log:  slog.Default().With("component", "service.preset"),
	}
}

func (s *Service) Create(ctx context.Context, p *preset.Preset) (int64, error) {
	const op = "service.preset.Create"
	log := s.log.With("op", op)

	if err := p.Validate(); err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		switch {
		case errors.Is(err, der.ErrConflict):
			return 0, preset.ErrPresetAlreadyExists
		case errors.Is(err, der.ErrNotFound):
			return 0, preset.ErrInvalidProductID
		}

		log.Error("repo.Create failed", slog.Any("err", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset created", slog.Int64("preset_id", id))
	return id, nil
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
		if errors.Is(err, der.ErrNotFound) {
			return preset.ErrPresetNotFound
		}
		log.Error("repo.Delete failed", slog.Int64("preset_id", id), slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
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
