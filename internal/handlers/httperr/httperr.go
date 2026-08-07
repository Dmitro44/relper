package httperr

import (
	"errors"
	"fmt"
	"net/http"

	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/danielgtaylor/huma/v2"
)

func Status(err error) int {
	if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
		switch appErr.Code {
		case apperr.ErrorCodeNotFound:
			return http.StatusNotFound
		case apperr.ErrorCodeForbidden:
			return http.StatusForbidden
		case apperr.ErrorCodeBadRequest:
			return http.StatusBadRequest
		case apperr.ErrorCodeUnauthorized:
			return http.StatusUnauthorized
		}
	}
	return http.StatusInternalServerError
}

func Map(err error) error {
	if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
		return huma.NewError(Status(appErr), fmt.Sprintf("%s: %s", appErr.Code, appErr.Message))
	}
	return err
}
