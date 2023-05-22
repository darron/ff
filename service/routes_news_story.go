package service

import (
	"errors"
	"net/http"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) GetNewsStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "that id does not exist")
	}
	ns, err := s.conf.NewsStoryRepository.Find(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if ns.ID == "" || ns.RecordID == "" {
		return c.JSON(http.StatusNotFound, "that id does not exist")
	}
	return c.JSON(http.StatusOK, ns)
}

func (s HTTPService) CreateNewsStory(c echo.Context) error {
	ns := &core.NewsStory{}
	if err := c.Bind(ns); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if ns.RecordID == "" {
		return c.JSON(http.StatusBadRequest, errors.New("must supply record ID with request"))
	}
	id, err := s.conf.NewsStoryRepository.Store(ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, id)
}
