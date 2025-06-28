package preset

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	der "github.com/Neimess/zorkin-store-project/internal/domain/error"
)

var (
	ErrPresetNotFound      = errors.New("preset not found")
	ErrPresetAlreadyExists = errors.New("preset already exists")
	ErrPresetInvalid       = errors.New("preset is invalid")
	ErrPresetItemsEmpty    = errors.New("preset items cannot be empty")
	ErrPresetItemsTooMany  = errors.New("preset items cannot exceed 100")
	ErrPresetItemsInvalid  = errors.New("preset items are invalid")
	ErrPresetNameTooLong   = errors.New("preset name is too long")
)

type PresetRepository interface {
	Create(ctx context.Context, p *domain.Preset) (int64, error)
	Get(ctx context.Context, id int64) (*domain.Preset, error)
	Delete(ctx context.Context, id int64) error
	ListDetailed(ctx context.Context) ([]domain.Preset, error)
	ListShort(ctx context.Context) ([]domain.Preset, error)
}

type Service struct {
	repo PresetRepository
	log  *slog.Logger
}

type Deps struct {
	repo PresetRepository
}

func NewDeps(repo PresetRepository) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("preset: missing PresetRepository")
	}
	return &Deps{repo: repo}, nil
}

func New(d *Deps) *Service {
	return &Service{
		repo: d.repo,
		log:  slog.Default().With("component", "service.preset"),
	}
}

func (ps *Service) Create(ctx context.Context, preset *domain.Preset) (int64, error) {
	const op = "service.preset.Create"
	log := ps.log.With("op", op)

	id, err := ps.repo.Create(ctx, preset)
	if err != nil {

		switch {
		case errors.Is(err, der.ErrConflict) || errors.Is(err, der.ErrValidation) || errors.Is(err, der.ErrBadRequest) || errors.Is(err, der.ErrNotFound):
			return 0, ErrPresetAlreadyExists
		default:
			log.Error("repo failed", slog.Any("error", err))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("preset created", slog.Int64("preset_id", id))
	return id, nil
}

func (ps *Service) Get(ctx context.Context, id int64) (*domain.Preset, error) {
	const op = "service.preset.Get"
	log := ps.log.With("op", op)

	preset, err := ps.repo.Get(ctx, id)
	if err != nil {

		if errors.Is(err, der.ErrNotFound) {
			return nil, der.ErrNotFound
		}
		log.Error("repo failed", slog.Int64("preset_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset retrieved", slog.Int64("preset_id", preset.ID))
	return preset, nil
}

func (ps *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.preset.Delete"
	log := ps.log.With("op", op)

	err := ps.repo.Delete(ctx, id)
	if err != nil {
		log.Error("repo failed", slog.Int64("preset_id", id), slog.Any("error", err))
		if errors.Is(err, der.ErrNotFound) {
			return der.ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset deleted", slog.Int64("preset_id", id))
	return nil
}

func (ps *Service) ListDetailed(ctx context.Context) ([]domain.Preset, error) {
	const op = "service.preset.List"
	log := ps.log.With("op", op)

	presets, err := ps.repo.ListDetailed(ctx)
	if err != nil {
		log.Error("repo failed", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset list retrieved", slog.Int("count", len(presets)))
	return presets, nil
}

func (ps *Service) ListShort(ctx context.Context) ([]domain.Preset, error) {
	const op = "service.preset.ListShort"
	log := ps.log.With("op", op)

	presets, err := ps.repo.ListShort(ctx)
	if err != nil {
		log.Error("repo failed", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset short list retrieved", slog.Int("count", len(presets)))
	return presets, nil
}
