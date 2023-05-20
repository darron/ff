package config

import (
	"flag"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/core"
)

type OptFunc func(*Opts)

type Opts struct {
	Debug               bool
	RecordRepository    core.RecordService
	NewsStoryRepository core.NewsStoryService
}

type App struct {
	Opts
}

func defaultOpts() Opts {
	return Opts{
		Debug: true,
	}
}

func withRedis(conn string) OptFunc {
	return func(opts *Opts) {
		opts.RecordRepository = redis.NewRecordRepository(conn)
		opts.NewsStoryRepository = redis.NewNewsStoryRepository(conn)
	}
}

func New() (*App, error) {
	var optFuncs []OptFunc
	var storage string
	var redis string

	flag.StringVar(&storage, "storage", "redis", "Storage for Data")
	flag.StringVar(&redis, "redis", "127.0.0.1:6379", "Redis Connection String")
	flag.Parse()

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
