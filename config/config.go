package config

import (
	"fmt"
	"os"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/core"
	"golang.org/x/exp/slog"
)

type OptFunc func(*Opts)

type Opts struct {
	JWTSecret           string
	Logger              *slog.Logger
	NewsStoryRepository core.NewsStoryService
	Port                string
	RecordRepository    core.RecordService
}

type App struct {
	Opts
}

var (
	defaultLogformat = "text"
	defaultLogLevel  = "debug"
	defaultPort      = "8080"
	defaultRedis     = "127.0.0.1:6379"
)

func defaultOpts() Opts {
	return Opts{
		Port: defaultPort,
	}
}

func WithRedis(conn string) OptFunc {
	return func(opts *Opts) {
		opts.RecordRepository = redis.NewRecordRepository(conn)
		opts.NewsStoryRepository = redis.NewNewsStoryRepository(conn)
	}
}

func WithJWTSecret(secret string) OptFunc {
	return func(opts *Opts) {
		opts.JWTSecret = secret
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

func New() (*App, error) {
	var optFuncs []OptFunc

	// Really basic default options without any configuration.
	// Moving configuration to cmd and will be calling `Get(WithOption())`
	optFuncs = append(optFuncs, WithLogger(defaultLogLevel, defaultLogformat))
	optFuncs = append(optFuncs, WithPort(defaultPort))
	optFuncs = append(optFuncs, WithRedis(defaultRedis))

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
	// TODO: Will need to update once we move to HTTPS
	// and a real domain name.
	return fmt.Sprintf("http://localhost:%s", a.Port)
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
