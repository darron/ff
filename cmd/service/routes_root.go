package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, "Nothing here")
}
