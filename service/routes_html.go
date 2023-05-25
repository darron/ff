package service

import (
	"net/http"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) Root(c echo.Context) error {
	// For display
	var dRecords []core.Record
	records, err := s.conf.RecordRepository.GetAll()
	// De-reference everything.
	for _, record := range records {
		dRecords = append(dRecords, *record)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "records", dRecords)
}

func (s HTTPService) IndividualRecord(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusNotFound, "id must not be blank")
	}
	r, err := s.conf.RecordRepository.Find(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "record", r)
}