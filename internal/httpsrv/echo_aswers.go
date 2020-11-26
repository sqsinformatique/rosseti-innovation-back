package httpsrv

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// BadRequest return err 400
func EchoBadRequest(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusBadRequest,
		BadRequest(err),
	)
}

// Unauthorized return err 401
func EchoUnauthorized(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusUnauthorized,
		Unauthorized(err),
	)
}

// Forbidden return err 403
func EchoForbidden(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusForbidden,
		Forbidden(err),
	)
}

// NotFound return err 404
func EchoNotFound(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusBadRequest,
		NotFound(err),
	)
}

// NotDeleted return err 406
func EchoNotDeleted(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusNotAcceptable,
		NotDeleted(err),
	)
}

// HasExpired return err 408
func EchoHasExpired(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusRequestTimeout,
		HasExpired(err),
	)
}

// CreateFailed return err 409
func EchoCreateFailed(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusConflict,
		CreateFailed(err),
	)
}

// NotUpdated return err 409
func EchoNotUpdated(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusConflict,
		NotUpdated(err),
	)
}

// InternalServerError return err 500
func EchoInternalServerError(ec echo.Context, err error) error {
	return ec.JSON(
		http.StatusInternalServerError,
		InternalServerError(err),
	)
}

func EchoOK(ec echo.Context, data interface{}) error {
	return ec.JSON(
		http.StatusOK,
		data,
	)
}

func EchoOkResult(ec echo.Context) error {
	return ec.JSON(
		http.StatusOK,
		OkResult(),
	)
}
