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

func TestRoot200(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
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
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, s.Root(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestIndividualRecord200(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/:id")
	c.SetParamNames("id")
	c.SetParamValues("asdf-1234-fdasa")
	if assert.NoError(t, s.IndividualRecord(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "OIC Impact")
	}
}

func TestIndividualRecord404(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/:id")
	c.SetParamNames("id")
	c.SetParamValues("")
	if assert.NoError(t, s.IndividualRecord(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestIndividualRecord500(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/:id")
	c.SetParamNames("id")
	c.SetParamValues("asdf-4312-fdass")
	if assert.NoError(t, s.IndividualRecord(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestGroup200(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/group/:group")
	c.SetParamNames("group")
	c.SetParamValues("mass")
	if assert.NoError(t, s.Group(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "OIC Impact")
	}
}

func TestGroup404(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/group/:group")
	c.SetParamNames("group")
	c.SetParamValues("")
	if assert.NoError(t, s.Group(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestGroup400(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepository{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/group/:group")
	c.SetParamNames("group")
	c.SetParamValues("not-a-group")
	if assert.NoError(t, s.Group(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestGroup500(t *testing.T) {
	s := HTTPService{}
	m := mockRecordRepositoryError{}
	conf, _ := config.Get()
	s.conf = conf
	s.conf.RecordRepository = m
	e := echo.New()
	tp, err := GetTemplates("../views/*.html")
	assert.NoError(t, err)
	e.Renderer = tp
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/records/group/:group")
	c.SetParamNames("group")
	c.SetParamValues("mass")
	if assert.NoError(t, s.Group(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestGetGroup(t *testing.T) {
	var records []*core.Record
	for i := 0; i < 200; i++ {
		r := core.FakeRecord()
		records = append(records, &r)
	}
	type test struct {
		input []*core.Record
		group string
		err   error
	}

	tests := []test{
		{input: records, group: "mass", err: nil},
		{input: records, group: "massfirearms", err: nil},
		{input: records, group: "massfirearmslicensed", err: nil},
		{input: records, group: "massother", err: nil},
		{input: records, group: "oic", err: nil},
		{input: records, group: "suicide", err: nil},
	}

	for _, test := range tests {
		_, got := GetGroup(test.group, test.input)
		if got != test.err {
			t.Errorf("Wanted: %q Got: %q", test.err, got)
		}
	}

	_, err := GetGroup("does-not-exist", records)
	assert.Error(t, err)
}
