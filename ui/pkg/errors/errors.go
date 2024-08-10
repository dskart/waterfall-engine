package common

import (
	"net/http"

	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func ErrorHTTPStatus(err appErrors.SanitizedError) int {
	switch err.(type) {
	case *appErrors.ResourceNotFoundError:
		return http.StatusNotFound
	case *appErrors.UserError:
		return http.StatusBadRequest
	case *appErrors.AuthorizationError:
		return http.StatusForbidden
	case *appErrors.AuthenticationError:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

type StatusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *StatusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
