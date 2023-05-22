package service

import (
	"encoding/json"
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

func TestGetNewsStory200(t *testing.T) {
	id := "asdf-2134-asdf-4321"
	s := HTTPService{}
	m := mockNewsStoryRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/stories/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	if assert.NoError(t, s.GetNewsStory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Let's look at the body too.
		theBody := rec.Body.String()
		// Does the body turn into a core.Record?
		_, err := core.UnmarshalJSONNewsStory(theBody)
		assert.NoError(t, err)
		// Grab ID from Body
		fromBody := gjson.Get(theBody, "id").String()
		assert.Equal(t, id, fromBody)
	}
}

func TestGetNewsStory404(t *testing.T) {
	s := HTTPService{}
	m := mockNewsStoryRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/stories/:id")
	c.SetParamNames("id")
	c.SetParamValues("")
	if assert.NoError(t, s.GetNewsStory(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestGetNewsStory500(t *testing.T) {
	id := "asdf-2134-asdf-4321"
	s := HTTPService{}
	m := mockNewsStoryRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/stories/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	if assert.NoError(t, s.GetNewsStory(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestCreateNewsStory200(t *testing.T) {
	s := HTTPService{}
	m := mockNewsStoryRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	recordJSON := core.FakeNewsStoryJSON()
	req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(recordJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateNewsStory(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestCreateNewsStory400(t *testing.T) {
	s := HTTPService{}
	m := mockNewsStoryRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader("not-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateNewsStory(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateNewsStory400MissingRecordID(t *testing.T) {
	s := HTTPService{}
	m := mockNewsStoryRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	ns := core.FakeNewsStory()
	// If this is missing - that's an error.
	ns.RecordID = ""
	recordJSON, _ := json.Marshal(ns)
	req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(string(recordJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateNewsStory(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateNewsStory500(t *testing.T) {
	s := HTTPService{}
	m := mockNewsStoryRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.NewsStoryRepository = m
	e := echo.New()
	recordJSON := core.FakeNewsStoryJSON()
	req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(recordJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.CreateNewsStory(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
