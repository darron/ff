package service

import (
	"net/http"

	"github.com/darron/ff/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPService struct {
	conf *config.App
}

func Get(conf *config.App) (*echo.Echo, error) {
	s := HTTPService{}

	s.conf = conf

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", root)
	e.GET("/record/:id", s.GetRecord)

	return e, nil
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Nothing here")
}

func (s HTTPService) GetRecord(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusNotFound, "That id does not exist")
	}
	r, err := s.conf.RecordRepository.Find(id)
	if err != nil {
		return c.String(http.StatusNotFound, "That id does not exist")
	}
	return c.JSON(http.StatusOK, r)
}
