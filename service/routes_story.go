package service

import (
	"net/http"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) CreateStory(c echo.Context) error {
	ns := &core.NewsStory{}
	if err := c.Bind(ns); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	id, err := s.conf.NewsStoryRepository.Store(ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, id)
}
