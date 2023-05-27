package service

import (
	"errors"
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

func (s HTTPService) Group(c echo.Context) error {
	group := c.Param("group")
	if group == "" {
		return c.String(http.StatusNotFound, "group must not be blank")
	}
	records, err := s.conf.RecordRepository.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// Let's get the subset
	newRecords, err := GetGroup(group, records)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Render(http.StatusOK, "records", newRecords)
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

func GetGroup(group string, records []*core.Record) ([]core.Record, error) {
	var finalGroup []core.Record
	var err error
	// This is pretty terrible - but it's just a test.
	switch group {
	case "mass":
		for _, record := range records {
			if record.Victims > 3 {
				finalGroup = append(finalGroup, *record)
			}
		}
	case "massfirearms":
		for _, record := range records {
			if record.Victims > 3 && record.Firearms.Bool {
				finalGroup = append(finalGroup, *record)
			}
		}
	case "massfirearmslicensed":
		for _, record := range records {
			if record.Victims > 3 && record.Firearms.Bool && record.Licensed.Bool {
				finalGroup = append(finalGroup, *record)
			}
		}
	case "massother":
		for _, record := range records {
			if record.Victims > 3 && !record.Firearms.Bool {
				finalGroup = append(finalGroup, *record)
			}
		}
	case "oic":
		for _, record := range records {
			if record.OICImpact.Bool {
				finalGroup = append(finalGroup, *record)
			}
		}
	case "suicide":
		for _, record := range records {
			if record.Suicide.Bool {
				finalGroup = append(finalGroup, *record)
			}
		}
	default:
		return finalGroup, errors.New("no group found")
	}
	return finalGroup, err
}
