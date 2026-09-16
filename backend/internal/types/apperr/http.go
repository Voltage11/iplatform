package apperr

import "net/http"

// Для не каст. ошибок возвращает 500 Internal Server Error
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	appErr, ok := AsAppError(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch appErr.Type {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
