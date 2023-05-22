package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/core"
	"golang.org/x/exp/slog"
)

type OptFunc func(*Opts)

type Opts struct {
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
	defaultStorage   = "redis"
)

func defaultOpts() Opts {
	return Opts{
		Port: defaultPort,
	}
}

func withRedis(conn string) OptFunc {
	return func(opts *Opts) {
		opts.RecordRepository = redis.NewRecordRepository(conn)
		opts.NewsStoryRepository = redis.NewNewsStoryRepository(conn)
	}
}

func withLogger(l *slog.Logger) OptFunc {
	return func(opts *Opts) {
		opts.Logger = l
	}
}

func setPort(port string) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

func New() (*App, error) {
	var logLevel string
	var logformat string
	var optFuncs []OptFunc
	var port string
	var redis string
	var storage string

	flag.StringVar(&logLevel, "loglevel", defaultLogLevel, "Default Log Level")
	flag.StringVar(&logformat, "logformat", defaultLogformat, "Log format: json or text")
	flag.StringVar(&port, "port", defaultPort, "HTTP Port")
	flag.StringVar(&redis, "redis", defaultRedis, "Redis Connection String")
	flag.StringVar(&storage, "storage", defaultStorage, "Storage for Data")

	flag.Parse()

	// Let's deal with logging immediately
	logger := GetLogger(logLevel, logformat)
	optFuncs = append(optFuncs, withLogger(logger))

	// If we've picked other options - add them to optFuncs
	if port != defaultPort {
		optFuncs = append(optFuncs, setPort(port))
	}
	if storage == "redis" && redis != "" {
		optFuncs = append(optFuncs, withRedis(redis))
	}

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
