package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

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
	c.SetParamValues("")
	if assert.NoError(t, s.GetRecord(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestGetRecord500(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepositoryError{}
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
	c.SetParamValues("does-not-matter")
	if assert.NoError(t, s.GetRecord(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestCreateRecord200(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	recordJSON := core.FakeRecordJSON()
	req := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(recordJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateRecord(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateRecord400(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader("not-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateRecord(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateRecord500(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	recordJSON := core.FakeRecordJSON()
	req := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(recordJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateRecord(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
