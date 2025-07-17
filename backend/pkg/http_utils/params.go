package http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func IDFromURL(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	if idStr == "" {
		return 0, fmt.Errorf("missing path parameter: %s", param)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid path parameter %s: %w", param, err)
	}
	return id, nil
}

func QueryInt64Param(r *http.Request, key string) (int64, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return 0, fmt.Errorf("query parameter %s is required", key)
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid query parameter %s: %w", key, err)
	}
	return value, nil
}
