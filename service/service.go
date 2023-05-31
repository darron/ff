package service

import (
	"html/template"
	"io"
	"time"

	cache "github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/memory"
	"github.com/darron/ff/config"
	"github.com/labstack/echo-contrib/prometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/guregu/null.v4"
)

type HTTPService struct {
	conf *config.App
}

var (
	APIPath                = "/api/v1"
	NewsStoriesAPIPath     = "/stories"
	NewsStoriesAPIPathFull = APIPath + NewsStoriesAPIPath
	RecordsAPIPath         = "/records"
	RecordsAPIPathFull     = APIPath + RecordsAPIPath

	// Middleware Cache settings
	cacheCapacity   = 10000
	cacheRefreshKey = "opn" // ?$cacheRefreshKey=true to a page to force a refresh
	cacheTTL        = 32 * time.Hour
	nonCachedPaths  = []string{"/api"}
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
	e.Use(middleware.RequestID())

	// Turn on JWT for the APIPath.
	// NOTE: This means unless you set JWT_SECRET deliberately
	// those endpoints are effectively locked by default.
	j := e.Group(APIPath)
	if s.conf.JWTSecret != "" {
		j.Use(echojwt.JWT([]byte(s.conf.JWTSecret)))
	}

	// Let's setup the in memory cache as middleware.
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cacheCapacity),
	)
	if err != nil {
		return e, err
	}
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(cacheTTL),
		cache.ClientWithRefreshKey(cacheRefreshKey),
		cache.ClientWithRestrictedPaths(nonCachedPaths),
	)
	if err != nil {
		return e, err
	}
	e.Use(cacheClient.Middleware())

	// Routes
	e.GET("/", s.Root)
	e.GET("/records/group/:group", s.Group)
	e.GET("/records/:id", s.IndividualRecord)

	// Infra
	e.GET("/version", s.Version)

	// This Group adds the APIPath to the full path when it
	// is created. This means we only need to supply what
	// is in addition to APIPath.
	// NOTE: When calling the endpoint in the cli, we have added
	// the "Full" option which gives the entire path including
	// APIPath again.
	j.GET(RecordsAPIPath+"/:id", s.GetRecord)
	j.GET(RecordsAPIPath, s.GetAllRecords)
	j.POST(RecordsAPIPath, s.CreateRecord)
	j.GET(NewsStoriesAPIPath+"/:id", s.GetNewsStory)
	j.POST(NewsStoriesAPIPath, s.CreateNewsStory)

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
