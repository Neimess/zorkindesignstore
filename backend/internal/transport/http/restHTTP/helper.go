package restHTTP

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
)

func idFromUrlToInt64(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id in URL param: %w", err)
	}
	return id, nil
}

func writeError(w http.ResponseWriter, statusCode int, msg string) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return easyjson.MarshalToWriter(dto.ErrorResponse{
		Message: msg,
	}, w)
}

func JSONCtx(ctx context.Context, w http.ResponseWriter, statusCode int, data easyjson.Marshaler, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := easyjson.MarshalToWriter(data, w); err != nil {
		logger.ErrorContext(ctx, "easyjson marshal failed", slog.String("error", err.Error()))
	}
}
