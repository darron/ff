package service

import (
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
	e.GET("/records/:id", s.GetRecord)
	e.GET("/records", s.GetAllRecords)
	e.POST("/records", s.CreateRecord)
	e.GET("/stories/:id", s.GetNewsStory)
	e.POST("/stories", s.CreateNewsStory)

	return e, nil
}
