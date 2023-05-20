package config

import (
	"flag"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/core"
)

type OptFunc func(*Opts)

type Opts struct {
	Debug               bool
	NewsStoryRepository core.NewsStoryService
	Port                string
	RecordRepository    core.RecordService
}

type App struct {
	Opts
}

var (
	defaultDebug   = true
	defaultPort    = "8080"
	defaultRedis   = "127.0.0.1:6379"
	defaultStorage = "redis"
)

func defaultOpts() Opts {
	return Opts{
		Debug: defaultDebug,
		Port:  defaultPort,
	}
}

func withRedis(conn string) OptFunc {
	return func(opts *Opts) {
		opts.RecordRepository = redis.NewRecordRepository(conn)
		opts.NewsStoryRepository = redis.NewNewsStoryRepository(conn)
	}
}

func setPort(port string) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

func setDebug(debug bool) OptFunc {
	return func(opts *Opts) {
		opts.Debug = debug
	}
}

func New() (*App, error) {
	var debug bool
	var optFuncs []OptFunc
	var port string
	var redis string
	var storage string

	flag.BoolVar(&debug, "debug", defaultDebug, "Debug logs")
	flag.StringVar(&port, "port", defaultPort, "HTTP Port")
	flag.StringVar(&redis, "redis", defaultRedis, "Redis Connection String")
	flag.StringVar(&storage, "storage", defaultStorage, "Storage for Data")

	flag.Parse()

	// If we've picked other options - add them to optFuncs
	if debug != defaultDebug {
		optFuncs = append(optFuncs, setDebug(debug))
	}
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
