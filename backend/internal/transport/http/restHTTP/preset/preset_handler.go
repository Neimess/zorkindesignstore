package preset

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	pSvc "github.com/Neimess/zorkin-store-project/internal/service/preset"
	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/Neimess/zorkin-store-project/pkg/validator"
)

type PresetService interface {
	Create(ctx context.Context, p *domain.Preset) (int64, error)
	Get(ctx context.Context, id int64) (*domain.Preset, error)
	Delete(ctx context.Context, id int64) error
	ListDetailed(ctx context.Context) ([]domain.Preset, error)
	ListShort(ctx context.Context) ([]domain.Preset, error)
}

type Deps struct {
	pSrv PresetService
}

func NewDeps(pSrv PresetService) (*Deps, error) {
	if pSrv == nil {
		return nil, errors.New("preset: missing PresetService")
	}
	return &Deps{
		pSrv: pSrv,
	}, nil
}

type Handler struct {
	srv PresetService
	log *slog.Logger
}

func New(d *Deps) *Handler {
	return &Handler{
		srv: d.pSrv,
		log: slog.Default().With("component", "transport.http.restHTTP.preset"),
	}
}

// Create godoc
// @Summary Create a new preset
// @Description Create a new preset with its items
// @Tags Preset
// @Accept json
// @Security BearerAuth
// @Param preset body dto.PresetRequest true "Preset data"// @Success 201 {object} dto.IDResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/admin/presets [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Create")
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	// Decode request
	var input dto.PresetRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteJSON(w, http.StatusBadRequest, "INVALID JSON")
		return
	}

	// Validate
	if err := validator.GetValidator().StructCtx(ctx, &input); err != nil {
		log.Warn("validation failed", slog.Any("error", err))
		httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid preset data")
		return
	}

	// Business logic
	preset := input.MapToPreset()
	id, err := h.srv.Create(ctx, preset)
	if err != nil {
		switch {
		case errors.Is(err, pSvc.ErrPresetAlreadyExists):
			httputils.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, pSvc.ErrPresetInvalid),
			errors.Is(err, pSvc.ErrPresetItemsEmpty),
			errors.Is(err, pSvc.ErrPresetItemsTooMany),
			errors.Is(err, pSvc.ErrPresetItemsInvalid),
			errors.Is(err, pSvc.ErrPresetNameTooLong):
			httputils.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			log.Error("create failed", slog.Any("error", err))
			httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	// Success
	httputils.WriteJSON(w, http.StatusCreated, dto.IDResponse{ID: id, Message: "Preset created successfully"})
}

// Get godoc
// @Summary Get preset by ID
// @Tags Preset
// @Produce json
// @Param id path int true "Preset ID"
// @Success 200 {object} dto.PresetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/presets/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Get")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	preset, err := h.srv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pSvc.ErrPresetNotFound) {
			httputils.WriteError(w, http.StatusNotFound, "preset not found")
			return
		}
		log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := dto.MapToDto(preset)
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete preset
// @Tags Preset
// @Security BearerAuth
// @Param id path int true "Preset ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/admin/presets/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Delete")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.srv.Delete(ctx, id)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, pSvc.ErrPresetNotFound):
		httputils.WriteError(w, http.StatusNotFound, "preset not found")
	default:
		log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}

// ListDetailed godoc
// @Summary List presets with items
// @Tags Preset
// @Produce json
// @Success 200 {array} dto.PresetResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/presets/detailed [get]
func (h *Handler) ListDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "ListDetailed")

	presets, err := h.srv.ListDetailed(ctx)
	if err != nil {
		log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make(dto.PresetListResponse, len(presets))
	for i, p := range presets {
		out[i] = *dto.MapToDto(&p)
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}

// ListShort godoc
// @Summary List presets basic info
// @Tags Preset
// @Produce json
// @Success 200 {array} dto.PresetShortResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/presets [get]
func (h *Handler) ListShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "ListShort")

	presets, err := h.srv.ListShort(ctx)
	if err != nil {
		log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make(dto.PresetShortListResponse, len(presets))
	for i, p := range presets {
		out[i] = dto.MapToShortDto(&p)
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}
