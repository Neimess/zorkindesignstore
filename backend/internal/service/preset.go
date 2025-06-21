package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
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

type PresetService struct {
	repo PresetRepository
	log  *slog.Logger
}

func NewPresetService(repo PresetRepository, log *slog.Logger) *PresetService {
	if repo == nil {
		panic("preset service: repo is nil")
	}
	return &PresetService{
		repo: repo,
		log:  logger.WithComponent(log, "service.preset"),
	}
}

func (ps *PresetService) Create(ctx context.Context, preset *domain.Preset) (int64, error) {
	const op = "service.preset.Create"
	log := ps.log.With("op", op)
	presetID, err := ps.repo.Create(ctx, preset)
	if err != nil {
		log.Error("repo failed", slog.Any("error", err))
		return 0, err
	}

	log.Info("preset created", slog.Int64("preset_id", presetID))
	return presetID, nil
}

func (ps *PresetService) Get(ctx context.Context, id int64) (*domain.Preset, error) {
	const op = "service.preset.Get"
	log := ps.log.With("op", op)

	preset, err := ps.repo.Get(ctx, id)
	switch {
	case errors.Is(err, repository.ErrPresetNotFound):
		return nil, ErrPresetNotFound
	case err != nil:
		log.Error("repo failed", slog.Int64("preset_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset retrieved", slog.Int64("preset_id", preset.ID))
	return preset, nil
}

func (ps *PresetService) Delete(ctx context.Context, id int64) error {
	const op = "service.preset.Delete"
	log := ps.log.With("op", op)

	if err := ps.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrPresetNotFound) {
			return ErrPresetNotFound
		}
		log.Error("repo failed", slog.Int64("preset_id", id), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("preset deleted", slog.Int64("preset_id", id))
	return nil
}

func (ps *PresetService) ListDetailed(ctx context.Context) ([]domain.Preset, error) {
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

func (ps *PresetService) ListShort(ctx context.Context) ([]domain.Preset, error) {
	// TODO обработка ошибок отсутсвия пресетов
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
