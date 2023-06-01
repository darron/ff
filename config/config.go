package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/core"
	"golang.org/x/exp/slog"
)

type OptFunc func(*Opts)

type Opts struct {
	JWTSecret           string
	JWTToken            string
	Logger              *slog.Logger
	MiddlewareHTMLCache bool
	NewsStoryRepository core.NewsStoryService
	Port                string
	RecordRepository    core.RecordService
	TLS                 TLS
}

type TLS struct {
	CacheDir    string
	DomainNames string
	Email       string
	Enable      bool
}

type App struct {
	Opts
}

var (
	defaultLogformat           = "text"
	defaultLogLevel            = "debug"
	defaultPort                = "8080"
	defaultRedis               = "127.0.0.1:6379"
	defaultHTMLMiddlewareCache = true
)

func defaultOpts() Opts {
	return Opts{
		Port:                defaultPort,
		MiddlewareHTMLCache: true,
	}
}

func WithRedis(conn string) OptFunc {
	return func(opts *Opts) {
		opts.RecordRepository = redis.RecordRepository{Conn: conn}
		opts.NewsStoryRepository = redis.NewsStoryRepository{Conn: conn}
	}
}

func WithMiddlewareHTMLCache(enabled bool) OptFunc {
	return func(opts *Opts) {
		opts.MiddlewareHTMLCache = enabled
	}
}

func WithJWTSecret(secret string) OptFunc {
	return func(opts *Opts) {
		opts.JWTSecret = secret
	}
}

func WithJWTToken(token string) OptFunc {
	return func(opts *Opts) {
		opts.JWTToken = token
	}
}

func WithLogger(level, format string) OptFunc {
	l := GetLogger(level, format)
	return func(opts *Opts) {
		opts.Logger = l
	}
}

func WithPort(port string) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

func WithTLS(tls TLS) OptFunc {
	return func(opts *Opts) {
		opts.TLS = tls
	}
}

func New() (*App, error) {
	var optFuncs []OptFunc

	// Really basic default options without any configuration.
	// Moving configuration to cmd and will be calling `Get(WithOption())`
	optFuncs = append(optFuncs, WithLogger(defaultLogLevel, defaultLogformat))
	optFuncs = append(optFuncs, WithPort(defaultPort))
	optFuncs = append(optFuncs, WithRedis(defaultRedis))
	optFuncs = append(optFuncs, WithMiddlewareHTMLCache(defaultHTMLMiddlewareCache))

	return Get(optFuncs...)
}

func Get(opts ...OptFunc) (*App, error) {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}
	app := App{
		Opts: o,
	}
	return &app, nil
}

func (a *App) GetHTTPEndpoint() string {
	protocol := "http"
	domain := "localhost"
	port := a.Port
	if a.TLS.DomainNames != "" {
		protocol = "https"
		domain = strings.Split(a.TLS.DomainNames, ",")[0]
		port = "443"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, domain, port)
}

func GetLogger(level, format string) *slog.Logger {
	var slogLevel slog.Level
	var slogHandler slog.Handler

	// Let's deal with level.
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	default:
		slogLevel = slog.LevelInfo
	}
	handlerOpts := slog.HandlerOptions{
		Level: slogLevel,
	}

	// Let's switch formats as desired.
	switch format {
	case "json":
		slogHandler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	default:
		slogHandler = slog.NewTextHandler(os.Stdout, &handlerOpts)
	}
	log := slog.New(slogHandler)

	return log
}

func (t TLS) Verify() error {
	if t.CacheDir == "" {
		return errors.New("Cache dir cannot be emtpy")
	}
	// Check to see if the cache dir exists - if it doesn't try to create it.
	if _, err := os.Open(t.CacheDir); os.IsNotExist(err) {
		// It doesn't exist - try to create it.
		err := os.MkdirAll(t.CacheDir, 0750)
		if err != nil {
			return err
		}
	}
	if t.DomainNames == "" {
		return errors.New("Domain names cannot be empty")
	}
	if t.Email == "" {
		return errors.New("Email address cannot be empty")
	}
	return nil
}
