package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type mockRecordRepository struct {
}

func (mrr mockRecordRepository) Find(id string) (*core.Record, error) {
	r := core.FakeRecord()
	r.ID = id
	return &r, nil
}

func (mrr mockRecordRepository) Store(r *core.Record) (string, error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return r.ID, nil
}

func TestGetRecord200(t *testing.T) {
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
		// Let's look at the body too.
		theBody := rec.Body.String()
		// Does the body turn into a core.Record?
		_, err := core.UnmarshalJSONRecord(theBody)
		assert.NoError(t, err)
		// Grab ID from Body
		fromBody := gjson.Get(theBody, "id").String()
		assert.Equal(t, id, fromBody)
	}
}

func TestGetRecord404(t *testing.T) {
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
	// We want to trigger a 404 here.
	c.SetParamValues("")
	if assert.NoError(t, s.GetRecord(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
