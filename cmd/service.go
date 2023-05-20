package cmd

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Debug bool
	Port  string
}

func Start(conf *Config) error {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", root)

	// Start server
	e.Logger.Fatal(e.Start(conf.Port))
	return nil
}

func root(c echo.Context) error {
	return c.Render(http.StatusOK, "root", "")
}
