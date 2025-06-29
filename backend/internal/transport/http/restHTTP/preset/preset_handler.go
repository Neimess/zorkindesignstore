package preset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/interfaces"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"
)

type PresetService interface {
	Create(ctx context.Context, p *preset.Preset) (int64, error)
	Get(ctx context.Context, id int64) (*preset.Preset, error)
	Delete(ctx context.Context, id int64) error
	ListDetailed(ctx context.Context) ([]preset.Preset, error)
	ListShort(ctx context.Context) ([]preset.Preset, error)
}

type Deps struct {
	pSrv      PresetService
	validator interfaces.Validator
	log       *slog.Logger
}

func NewDeps(validator interfaces.Validator, log *slog.Logger, pSrv PresetService) (Deps, error) {
	if pSrv == nil {
		return Deps{}, errors.New("preset: missing PresetService")
	}
	if validator == nil {
		return Deps{}, errors.New("preset: missing validator")
	}
	if log == nil {
		return Deps{}, errors.New("preset: missing logger")
	}
	return Deps{
		pSrv:      pSrv,
		validator: validator,
		log:       log.With("component", "restHTTP.preset"),
	}, nil
}

type Handler struct {
	srv PresetService
	log *slog.Logger
	val interfaces.Validator
}

func New(d Deps) *Handler {
	return &Handler{
		srv: d.pSrv,
		log: d.log,
		val: d.validator,
	}
}

// Create godoc
// @Summary Create a new preset
// @Description Create a new preset with its items
// @Tags Preset
// @Accept json
// @Security BearerAuth
// @Param preset body dto.PresetRequest true "Preset data"// @Success 201 {object} dto.IDResponse
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 409 {object} httputils.ErrorResponse
// @Failure 422 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /api/admin/presets [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Create")
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	// 1) Decode
	var in dto.PresetRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %s", err.Error()))
		return
	}

	// 2) Validate struct tags
	if err := h.val.StructCtx(ctx, &in); err != nil {
		log.Warn("validation failed", slog.Any("error", err))

		var verrs = make([]dto.FieldError, 0)
		for _, fe := range err.(validator.ValidationErrors) {
			verrs = append(verrs, dto.FieldError{
				Field: fe.Field(),
				Tag:   fe.Tag(),
				Value: fe.Param(),
			})
		}
		httputils.WriteJSON(w, http.StatusUnprocessableEntity, dto.ValidationErrorResponse{
			Message: "Validation failed",
			Errors:  verrs,
		})
		return
	}

	// 3) Business logic
	p := in.MapToPreset()
	id, err := h.srv.Create(ctx, p)
	if err != nil {
		switch {
		case errors.Is(err, preset.ErrPresetAlreadyExists):
			httputils.WriteError(w, http.StatusConflict, "A preset with this name already exists")
		case errors.Is(err, preset.ErrNoItems):
			httputils.WriteError(w, http.StatusUnprocessableEntity, "At least one item is required")
		case errors.Is(err, preset.ErrNameTooLong):
			httputils.WriteError(w, http.StatusUnprocessableEntity, "Name must be at most 100 characters")
		case errors.Is(err, preset.ErrTotalPriceMismatch):
			httputils.WriteError(w, http.StatusUnprocessableEntity, "Total price must equal sum of item prices")
		case errors.Is(err, preset.ErrInvalidProductID):
			httputils.WriteError(w, http.StatusUnprocessableEntity, "One or more product IDs are invalid")
		default:
			log.Error("create failed", slog.Any("error", err))
			httputils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/admin/presets/%d", id))
	httputils.WriteJSON(w, http.StatusCreated, dto.PresetResponseID{PresetID: id, Message: "Preset created successfully"})
}

// Get godoc
// @Summary Get preset by ID
// @Tags Preset
// @Produce json
// @Param id path int true "Preset ID"
// @Success 200 {object} dto.PresetResponse
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /api/presets/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Get")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	p, err := h.srv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, preset.ErrPresetNotFound) {
			httputils.WriteError(w, http.StatusNotFound, "preset not found")
			return
		}
		log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := dto.MapDomainToDto(p)
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete preset
// @Tags Preset
// @Security BearerAuth
// @Param id path int true "Preset ID"
// @Success 204
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
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
	case errors.Is(err, preset.ErrPresetNotFound):
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
// @Failure 500 {object} httputils.ErrorResponse
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

	out := make([]dto.PresetResponse, len(presets))
	for i, p := range presets {
		out[i] = *dto.MapDomainToDto(&p)
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}

// ListShort godoc
// @Summary List presets basic info
// @Tags Preset
// @Produce json
// @Success 200 {array} dto.PresetShortResponse
// @Failure 500 {object} httputils.ErrorResponse
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

	out := make([]dto.PresetShortResponse, len(presets))
	for i, p := range presets {
		out[i] = dto.MapDomainToShortDTO(&p)
	}
	httputils.WriteJSON(w, http.StatusOK, out)
}
