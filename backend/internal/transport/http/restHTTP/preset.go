package restHTTP

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/service"
	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
	"github.com/Neimess/zorkin-store-project/pkg/validator"
	"github.com/mailru/easyjson"
)

type PresetService interface {
	Create(ctx context.Context, p *domain.Preset) (int64, error)
	Get(ctx context.Context, id int64) (*domain.Preset, error)
	Delete(ctx context.Context, id int64) error
	ListDetailed(ctx context.Context) ([]domain.Preset, error)
	ListShort(ctx context.Context) ([]domain.Preset, error)
}

type PresetHandler struct {
	srv PresetService
	log *slog.Logger
}

func NewPresetHandler(service PresetService, log *slog.Logger) *PresetHandler {
	if service == nil {
		panic("preset handler: service is nil")
	}
	return &PresetHandler{
		srv: service,
		log: logger.WithComponent(log, "transport.http.restHTTP.preset"),
	}
}

// CreatePreset 	godoc
// @Summary 		Create a new preset
// @Description 	Create a new preset with its items
// @Tags 			Preset
// @Accept 			json
// @Security  		BearerAuth
// @Param 			preset body dto.PresetRequest true "Preset data"
// @Success 		201 {object} dto.IDResponse "Preset created successfully"
// @Failure 		400 {object} dto.ErrorResponse "Invalid JSON or validation error"
// @Failure 		422 {object} dto.ErrorResponse "Unprocessable Entity"
// @Failure 		500 {object} dto.ErrorResponse "Internal Server Error"
// @Router 			/api/admin/presets [post]
func (h *PresetHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.Create")
	validator := validator.GetValidator()
	defer r.Body.Close()

	var input dto.PresetRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &input); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := validator.StructCtx(ctx, &input); err != nil {
		log.Warn("validation failed", slog.Any("error", err))
		httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid product data")
		return
	}

	preset := input.MapToPreset()

	presetID, err := h.srv.Create(ctx, preset)
	if err != nil {
		h.log.Error("create preset", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, dto.IDResponse{
		ID:      presetID,
		Message: "Preset created successfully",
	})
}

// GetPreset godoc
// @Summary      Get preset by ID
// @Tags         Preset
// @Produce      json
// @Param        id   path      int  true  "Preset ID"
// @Success      200 {object} dto.PresetResponse
// @Failure      400 {object} dto.ErrorResponse "invalid id"
// @Failure      404 {object} dto.ErrorResponse "not found"
// @Failure      500 {object} dto.ErrorResponse "internal error"
// @Router       /api/presets/{id} [get]
func (h *PresetHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "preset.Get")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	preset, err := h.srv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrPresetNotFound) {
			httputils.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := dto.MapToDto(preset)
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// DeletePreset godoc
// @Summary      Delete preset
// @Tags         Preset
// @Security     BearerAuth
// @Param        id   path      int  true  "Preset ID"
// @Success      204
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/admin/presets/{id} [delete]
func (h *PresetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "preset.Delete")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch err := h.srv.Delete(ctx, id); {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, service.ErrPresetNotFound):
		httputils.WriteError(w, http.StatusNotFound, "not found")
	default:
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}

// ListDetailed godoc
// @Summary      List presets with it's items
// @Description  Get a list of all presets with their items
// @Tags         Preset
// @Produce      json
// @Success      200 {object} dto.PresetListResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/presets/detailed [get]
func (h *PresetHandler) ListDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "preset.ListDetailed")

	presets, err := h.srv.ListDetailed(ctx)
	if err != nil {
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make(dto.PresetListResponse, len(presets))
	for i, p := range presets {
		out[i] = *dto.MapToDto(&p)
	}

	httputils.WriteJSON(w, http.StatusOK, out)
}

// ListShortPreset godoc
// @Summary      List presets with basic info
// @Description  Get a list of all presets with basic info (short version)
// @Tags         Preset
// @Produce      json
// @Success      200 {object} dto.PresetShortListResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/presets [get]
func (h *PresetHandler) ListShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "preset.List")

	presets, err := h.srv.ListShort(ctx)
	if err != nil {
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make(dto.PresetShortListResponse, len(presets))
	for i, p := range presets {
		out[i] = dto.MapToShortDto(&p)
	}

	httputils.WriteJSON(w, http.StatusOK, out)
}