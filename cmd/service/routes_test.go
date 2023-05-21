package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockRecordRepository struct {
}

func (mrr mockRecordRepository) Find(id string) (*core.Record, error) {
	r := core.FakeRecord()
	r.ID = id
	return &r, nil
}

func (mrr mockRecordRepository) Store(r *core.Record) (string, error) {
	return r.ID, nil
}

func TestGetRecord(t *testing.T) {
	id := "asdf-2134-asdf-4321"
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	if assert.NoError(t, s.GetRecord(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, id, rec.Body.String())
	}
}
