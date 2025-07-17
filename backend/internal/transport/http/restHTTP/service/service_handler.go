package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/service/dto"
	http_utils "github.com/Neimess/zorkin-store-project/pkg/http_utils"
)

type ServiceService interface {
	Create(ctx context.Context, s *domService.Service) (*domService.Service, error)
	Get(ctx context.Context, id int64) (*domService.Service, error)
	List(ctx context.Context) ([]domService.Service, error)
	Update(ctx context.Context, s *domService.Service) (*domService.Service, error)
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	Log *slog.Logger
	Srv ServiceService
}

func NewDeps(log *slog.Logger, srv ServiceService) (Deps, error) {
	if srv == nil {
		return Deps{}, errors.New("service: missing service")
	}
	if log == nil {
		return Deps{}, errors.New("service: missing logger")
	}
	return Deps{Log: log.With("component", "restHTTP.service"), Srv: srv}, nil
}

type Handler struct {
	srv ServiceService
	log *slog.Logger
}

func New(d Deps) *Handler {
	return &Handler{srv: d.Srv, log: d.Log}
}

// Create godoc
// @Summary      Create service
// @Description  Создать услугу
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        data body dto.ServiceRequest true "Service data"
// @Success      201 {object} dto.ServiceResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      409 {object} http_utils.ErrorResponse
// @Failure      422 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/services [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "Create")

	req, ok := http_utils.DecodeAndValidate[dto.ServiceRequest](w, r, log)
	if !ok {
		return
	}
	service := dto.MapToDomain(req)
	created, err := h.srv.Create(ctx, service)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(created)
	w.Header().Set("Location", fmt.Sprintf("/api/services/%d", resp.ID))
	http_utils.WriteJSON(w, http.StatusCreated, resp)
}

// PublicGet godoc
// @Summary      Get service by ID
// @Description  Получить услугу по ID
// @Tags         services
// @Produce      json
// @Param        id path int true "Service ID"
// @Success      200 {object} dto.ServiceResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      404 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/services/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	service, err := h.srv.Get(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(service)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// PublicList godoc
// @Summary      List services (public)
// @Description  Получить публичный список услуг
// @Tags         services
// @Produce      json
// @Success      200 {array} dto.ServiceResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/services [get]
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
// @Summary      Update service
// @Description  Обновить услугу по ID
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Service ID"
// @Param        data body dto.ServiceRequest true "Service data"
// @Success      200 {object} dto.ServiceResponse
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      409 {object} http_utils.ErrorResponse
// @Failure      422 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/services/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	req, ok := http_utils.DecodeAndValidate[dto.ServiceRequest](w, r, h.log)
	if !ok {
		return
	}
	service := dto.MapToDomain(req)
	service.ID = id
	updated, err := h.srv.Update(ctx, service)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapToResponse(updated)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete service
// @Description  Удалить услугу по ID
// @Tags         services
// @Security     BearerAuth
// @Param        id path int true "Service ID"
// @Success      204
// @Failure      400 {object} http_utils.ErrorResponse
// @Failure      404 {object} http_utils.ErrorResponse
// @Failure      500 {object} http_utils.ErrorResponse
// @Router       /api/admin/services/{id} [delete]
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
	case errors.Is(err, domService.ErrServiceNotFound):
		http_utils.WriteError(w, http.StatusNotFound, "service not found")
	case errors.Is(err, domService.ErrServiceAlreadyExists):
		http_utils.WriteError(w, http.StatusConflict, "service already exists")
	case errors.Is(err, domService.ErrEmptyName):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "name must not be empty")
	case errors.Is(err, domService.ErrNameTooLong):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "name is too long")
	default:
		h.log.Error("service error", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
