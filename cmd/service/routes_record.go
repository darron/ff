package service

import (
	"net/http"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) GetRecord(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "That id does not exist")
	}
	r, err := s.conf.RecordRepository.Find(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, r)
}

func (s HTTPService) CreateRecord(c echo.Context) error {
	r := &core.Record{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	id, err := s.conf.RecordRepository.Store(r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, id)
}
