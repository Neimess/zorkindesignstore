package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/auth/dto"
	http_utils "github.com/Neimess/zorkin-store-project/pkg/http_utils"
)

type AuthService interface {
	GenerateToken(userID string) (string, error)
}

type Deps struct {
	srv AuthService
	log *slog.Logger
}

func NewDeps(log *slog.Logger, srv AuthService) (Deps, error) {
	if srv == nil {
		return Deps{}, fmt.Errorf("auth: missing AuthService")
	}
	if log == nil {
		return Deps{}, fmt.Errorf("auth: missing Logger")
	}
	return Deps{
		srv: srv,
		log: log.With("component", "restHTTP.auth"),
	}, nil
}

type Handler struct {
	srv AuthService
	log *slog.Logger
}

func New(d Deps) *Handler {

	return &Handler{
		srv: d.srv,
		log: d.log,
	}
}

// Login godoc
// @Summary      Login user
// @Description  Generates a Bearer token for the user
// @Tags         auth
// @Produce      json
// @Param        secret_admin_key  path      string  true  "Secret admin key for login, injected via route"
// @Success      201  {object}  dto.TokenResponse  "Returns generated token"
// @Failure      401  {object}  http_utils.ErrorResponse  "Unauthorized access"
// @Failure      500  {object}  http_utils.ErrorResponse  "Internal server error"
// @Router       /api/admin/auth/{secret_admin_key} [get]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.CreateToken"
	log := h.log.With("op", op)
	userID := "ADMIN"

	token, err := h.srv.GenerateToken(userID)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	http_utils.WriteJSON(w, http.StatusCreated, dto.TokenResponse{
		Token: token,
	})
}
