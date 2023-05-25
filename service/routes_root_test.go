package service

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darron/ff/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoot200(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp := &Template{
		templates: template.Must(template.ParseGlob("../views/*.html")),
	}
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.Root(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "OIC Impact")
	}
}

func TestRoot500(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp := &Template{
		templates: template.Must(template.ParseGlob("../views/*.html")),
	}
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.Root(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
