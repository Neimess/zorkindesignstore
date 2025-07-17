package coefficients

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/coefficients/dto"
	http_utils "github.com/Neimess/zorkin-store-project/pkg/http_utils"
)

type CoefficientService interface {
	Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error)
	Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error)
	List(ctx context.Context) ([]domCoeff.Coefficient, error)
	Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error)
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	Log *slog.Logger
	Srv CoefficientService
}

func NewDeps(log *slog.Logger, srv CoefficientService) (Deps, error) {
	if srv == nil {
		return Deps{}, errors.New("coefficients: missing service")
	}
	if log == nil {
		return Deps{}, errors.New("coefficients: missing logger")
	}
	return Deps{Log: log.With("component", "restHTTP.coefficients"), Srv: srv}, nil
}

type Handler struct {
	srv CoefficientService
	log *slog.Logger
}

func New(d Deps) *Handler {
	return &Handler{srv: d.Srv, log: d.Log}
}

// Create godoc
// @Summary      Create coefficient
// @Description  Создать коэффициент
// @Tags         coefficients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body dto.CoefficientRequest true "Coefficient data"
// @Success      201 {object} dto.CoefficientResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      409 {object} http_utils.ErrorResponse
// @Failure      422 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/coefficients [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Create")

	req, ok := http_utils.DecodeAndValidate[dto.CoefficientRequest](w, r, log)
	if !ok {
		return
	}
	coeff := dto.MapToDomain(req)
	created, err := h.srv.Create(ctx, coeff)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(created)
	w.Header().Set("Location", fmt.Sprintf("/api/admin/coefficients/%d", resp.ID))
	http_utils.WriteJSON(w, http.StatusCreated, resp)
}

// Get godoc
// @Summary      Get coefficient by ID
// @Description  Получить коэффициент по ID
// @Tags         coefficients
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Coefficient ID"
// @Success      200 {object} dto.CoefficientResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      404 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/coefficients/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	coeff, err := h.srv.Get(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(coeff)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// List godoc
// @Summary      List coefficients
// @Description  Получить список коэффициентов
// @Tags         coefficients
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.CoefficientResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/coefficients [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := h.srv.List(ctx)
	if err != nil {
		h.log.Error("service error", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	resp := dto.MapToResponseList(list)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary      Update coefficient
// @Description  Обновить коэффициент по ID
// @Tags         coefficients
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Coefficient ID"
// @Param        data body dto.CoefficientRequest true "Coefficient data"
// @Success      200 {object} dto.CoefficientResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      409 {object} http_utils.ErrorResponse
// @Failure      422 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/coefficients/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	req, ok := http_utils.DecodeAndValidate[dto.CoefficientRequest](w, r, h.log)
	if !ok {
		return
	}
	coeff := dto.MapToDomain(req)
	coeff.ID = id
	updated, err := h.srv.Update(ctx, coeff)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(updated)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete coefficient
// @Description  Удалить коэффициент по ID
// @Tags         coefficients
// @Security     BearerAuth
// @Param        id path int true "Coefficient ID"
// @Success      204
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      404 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/coefficients/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.srv.Delete(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domCoeff.ErrCoefficientNotFound):
		http_utils.WriteError(w, http.StatusNotFound, "coefficient not found")
	case errors.Is(err, domCoeff.ErrCoefficientAlreadyExists):
		http_utils.WriteError(w, http.StatusConflict, "coefficient already exists")
	case errors.Is(err, domCoeff.ErrEmptyName):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "name must not be empty")
	case errors.Is(err, domCoeff.ErrNameTooLong):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "name is too long")
	default:
		h.log.Error("service error", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
