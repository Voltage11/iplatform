package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Voltage11/iplatform/internal/types/apperr"
	"github.com/Voltage11/iplatform/internal/types/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// WriteJSON отправляет JSON-ответ с заданным статусом
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteErrorString отправляет ошибку в формате JSON
func WriteErrorString(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, map[string]string{"error": message})
}

// WriteError отправляет ошибку в формате JSON
func WriteError(w http.ResponseWriter, err error) {
	var code int
	var message string

	switch {
	case err == nil:
		code = http.StatusInternalServerError
		message = "неизвестная ошибка"
	default:
		var validationErrs validator.ValidationErrors
		var appErr *apperr.AppError
		switch {
		case errors.As(err, &validationErrs):
			code = http.StatusBadRequest
			message = validationErrs.Error()
		case errors.As(err, &appErr):
			code = apperr.HTTPStatusFromError(appErr)
			message = appErr.Message
		default:
			code = http.StatusInternalServerError
			message = "внутренняя ошибка сервера"
		}
	}

	WriteJSON(w, code, map[string]string{"error": message})
}

// ParseUUID извлекает и валидирует UUID из параметра URL
func ParseUUID(r *http.Request, paramName string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, paramName)
	return uuid.Parse(idStr)
}

// ParsePagination извлекает page и limit из query-параметров
func ParsePagination(r *http.Request) pagination.PaginationRequest {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return pagination.NewPaginationRequest(page, limit)
}
