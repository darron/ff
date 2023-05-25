package service

import (
	"html/template"
	"io"

	"github.com/darron/ff/config"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/guregu/null.v4"
)

type HTTPService struct {
	conf *config.App
}

var (
	APIPath            = "/api/v1"
	NewsStoriesAPIPath = APIPath + "/stories"
	RecordsAPIPath     = APIPath + "/records"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Get(conf *config.App) (*echo.Echo, error) {
	s := HTTPService{}

	s.conf = conf

	// Echo instance
	e := echo.New()

	// Enable Prometheus
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Let's allow some static files
	e.Static("/", "public")

	// Let's setup the templates.
	t, err := GetTemplates("views/*.html")
	if err != nil {
		return nil, err
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", s.Root)
	e.GET("/records/:id", s.IndividualRecord)
	e.GET(RecordsAPIPath+"/:id", s.GetRecord)
	e.GET(RecordsAPIPath, s.GetAllRecords)
	e.POST(RecordsAPIPath, s.CreateRecord)
	e.GET(NewsStoriesAPIPath+"/:id", s.GetNewsStory)
	e.POST(NewsStoriesAPIPath, s.CreateNewsStory)

	return e, nil
}

func nullbool(n null.Bool) string {
	if n.Valid && n.Bool {
		return "Yes"
	} else if n.Valid && !n.Bool {
		return "No"
	}
	return ""
}

func GetTemplates(p string) (*Template, error) {
	// Let's setup templates with custom funcs.
	t := &Template{
		templates: template.New(""),
	}
	funcMap := template.FuncMap{
		"nullbool": nullbool,
	}
	t.templates = t.templates.Funcs(funcMap)
	_, err := t.templates.ParseGlob(p)
	if err != nil {
		return nil, err
	}
	return t, nil
}
