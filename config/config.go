package config

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/darron/ff/adaptors/redis"
	"github.com/darron/ff/adaptors/sqlite"
	"github.com/darron/ff/config/migrations"
	"github.com/darron/ff/core"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/exp/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
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
	CertFile    string
	DomainNames string
	Email       string
	Enable      bool
	KeyFile     string
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

func WithSQLite(file string) OptFunc {
	// Does file exist? If it doesn't - create it.
	if _, err := os.Stat(file); err != nil {
		err := createSQLite3Database(file)
		if err != nil {
			log.Fatal(err)
		}
		// If you've created it - migrate it.
		err = migrateSQLite3Database(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	return func(opts *Opts) {
		opts.RecordRepository = sqlite.RecordRepository{Filename: file}
		opts.NewsStoryRepository = sqlite.NewsStoryRepository{Filename: file}
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

func (t TLS) LetsEncryptVerify() error {
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

func (t TLS) StaticCredentialsVerify() error {
	_, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return err
	}
	return nil
}

func (t TLS) StaticCredentialsTLSConfig() (*tls.Config, error) {
	var tlsConfig *tls.Config
	cer, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return tlsConfig, err
	}
	tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	return tlsConfig, nil
}

func (t TLS) LetsEncryptTLSConfig() *tls.Config {
	domains := strings.Split(t.DomainNames, ",")
	autoTLSManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(t.CacheDir),
		Email:      t.Email,
		HostPolicy: autocert.HostWhitelist(domains...),
	}
	tlsConfig := tls.Config{
		GetCertificate: autoTLSManager.GetCertificate,
		NextProtos:     []string{acme.ALPNProto},
	}
	return &tlsConfig
}

func createSQLite3Database(file string) error {
	log.Println("Creating SQLite3 db.")
	newDB, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	db := sqlx.NewDb(newDB, "sqlite3")
	_, err = db.Exec("PRAGMA foreign_keys=OFF;")
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

func migrateSQLite3Database(file string) error {
	log.Println("Migrating SQLite3 db.")
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", d, "sqlite3://"+file)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
