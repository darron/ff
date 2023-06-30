package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s HTTPService) RecordForm(c echo.Context) error {
	// TODO: Going to have to figure out how to protect this form.
	return c.Render(http.StatusOK, "new-record", nil)
}
