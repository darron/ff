package service

import (
	"html/template"
	"io"

	"github.com/darron/ff/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// For when we want to use templates.
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
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
