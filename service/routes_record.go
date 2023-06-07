package service

import (
	"net/http"

	"github.com/darron/ff/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func checkJWTToken(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusBadRequest, "JWT token missing or invalid")
	}
	_, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, "failed to cast claims as jwt.MapClaims")
	}
	return nil
}

func (s HTTPService) GetRecord(c echo.Context) error {
	err := checkJWTToken(c)
	if err != nil {
		return err
	}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "id must not be blank")
	}
	r, err := s.conf.RecordRepository.Find(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if r.ID == "" {
		return c.JSON(http.StatusNotFound, "that id does not exist")
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

func (s HTTPService) GetAllRecords(c echo.Context) error {
	// Get all records.
	records, err := s.conf.RecordRepository.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, records)
}
