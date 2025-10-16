package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// serviceError is an interface to check for custom service errors
// that contain an HTTP status code.
type serviceError interface {
	GetCode() int
}

// HTTPErrorHandler is a custom error handler for Echo.
func HTTPErrorHandler(err error, c echo.Context) {
	var se serviceError
	if errors.As(err, &se) {
		c.JSON(se.GetCode(), map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
