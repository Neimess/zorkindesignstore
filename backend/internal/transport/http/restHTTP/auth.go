package restHTTP

import (
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

type AuthService interface {
	GenerateToken(userID string) (string, error)
}

type AuthHandler struct {
	srv AuthService
	log *slog.Logger
}

func NewAuthHandler(srv AuthService, log *slog.Logger) *AuthHandler {
	if srv == nil {
		panic("auth handler: AuthService is nil")
	}
	return &AuthHandler{
		srv: srv,
		log: logger.WithComponent(slog.Default(), "handler.auth"),
	}
}

// Login godoc
// @Summary      Login user
// @Description  Generates a Bearer token for the user
// @Tags         auth
// @Produce      json
// @Param        secret_admin_key  path      string  true  "Secret admin key for login, injected via route"
// @Success      201  {object}  dto.TokenResponse  "Returns generated token"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized access"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /api/auth/login/{secret_admin_key} [get]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.CreateToken"
	log := h.log.With("op", op)
	userID := "ADMIN"

	token, err := h.srv.GenerateToken(userID)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusCreated, dto.TokenResponse{
		Token: token,
	})
}
