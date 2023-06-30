package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/go-faker/faker/v4"
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "HTTP Service Commands",
		Run: func(cmd *cobra.Command, args []string) {
			StartService()
		},
	}

	defaultStorageLayer = "sqlite3"
	storageLayer        string

	defaultRedisConn = "127.0.0.1:6379"
	redisConn        string

	defaultSQLiteFile = "ff.db"
	sqliteFile        string

	jwtSecret string

	defaultProfilingEnabled = false
	profilingEnabled        bool

	defaultMiddlewareHTMLCacheEnabled = true
	middlewareHTMLCacheEnabled        bool
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&storageLayer, "storage", "", GetENVVariable("STORAGE", defaultStorageLayer), "Storage Layer: redis or sqlite3")
	serviceCmd.Flags().StringVarP(&sqliteFile, "sqlite3", "", GetENVVariable("SQLITE3", defaultSQLiteFile), "SQLite3 Filename")
	serviceCmd.Flags().StringVarP(&redisConn, "redisConn", "r", GetENVVariable("REDIS", defaultRedisConn), "Redis connection string")
	serviceCmd.Flags().StringVarP(&jwtSecret, "jwtSecret", "", GetENVVariable("JWT_SECRET", defaultJWTSecret()), "JWT Secret")
	serviceCmd.Flags().BoolVarP(&profilingEnabled, "profiling", "", GetBoolENVVariable("PROFILING_ENABLED", defaultProfilingEnabled), "Enable Datadog tracing and profiling")
	serviceCmd.Flags().BoolVarP(&middlewareHTMLCacheEnabled, "htmlcache", "", GetBoolENVVariable("HTMLCACHE_ENABLED", defaultMiddlewareHTMLCacheEnabled), "Enable Middleware Cache")
}

func StartService() {
	// Setup some options
	var opts []config.OptFunc
	var tlsConfig *tls.Config

	opts = append(opts, config.WithPort(port))
	opts = append(opts, config.WithLogger(logLevel, logFormat))
	opts = append(opts, config.WithMiddlewareHTMLCache(middlewareHTMLCacheEnabled))

	// Pick the storage layer and do the things.
	switch storageLayer {
	case "redis":
		log.Println("Enabling Redis storage")
		opts = append(opts, config.WithRedis(redisConn))
	case "sqlite3":
		log.Println("Enabling SQLite storage")
		opts = append(opts, config.WithSQLite(sqliteFile))
	default:
		log.Fatal("Must pick supported storage layer.")
	}

	// Let's enable JWT if it's defined.
	if jwtSecret != "" {
		log.Println("Enabling JWT")
		opts = append(opts, config.WithJWTSecret(jwtSecret))
	}

	// Let's turn on TLS with LetsEncrypt
	// Setup the config here.
	if enableTLS && enableTLSLetsEncrypt {
		log.Println("Enabling LetsEncrypt")
		tlsVar := config.TLS{
			CacheDir:    tlsCache,
			DomainNames: tlsDomains,
			Email:       tlsEmail,
			Enable:      enableTLS,
		}
		err := tlsVar.LetsEncryptVerify()
		if err != nil {
			log.Fatal(err)
		}
		// Let's setup the service http.Server tls.Config
		tlsConfig = tlsVar.LetsEncryptTLSConfig()
		opts = append(opts, config.WithTLS(tlsVar))
	}

	// If we have manually generated certs - let's use those for HTTPS
	// Setup the config here.
	if enableTLS && !enableTLSLetsEncrypt && (tlsCert != "") && (tlsKey != "") {
		log.Println("Enabling TLS with manual certs")
		tlsVar := config.TLS{
			CertFile:    tlsCert,
			DomainNames: tlsDomains,
			Enable:      enableTLS,
			KeyFile:     tlsKey,
		}
		err := tlsVar.StaticCredentialsVerify()
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig, err = tlsVar.StaticCredentialsTLSConfig()
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, config.WithTLS(tlsVar))
	}

	// Let's get the config for the app
	conf, err := config.Get(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Let's setup DD tracing and profiling:
	// NOTE: Make sure to set DD_ENV
	// for each place you're running this.
	if profilingEnabled {
		log.Println("Turning on Datadog Tracing")
		tracer.Start(
			tracer.WithService("ff"),
		)
		defer tracer.Stop()

		err = profiler.Start(
			profiler.WithService("ff"),
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer profiler.Stop()
	}

	conf.Logger.Info("Starting HTTP Service")
	s, err := service.Get(conf, "views/*.html")
	if err != nil {
		conf.Logger.Error(err.Error())
		os.Exit(1)
	}

	if logLevel == "debug" {
		for _, route := range s.Routes() {
			fmt.Printf("Route: %s %s\n", route.Method, route.Path)
		}
	}

	// If we are going to turn on TLS - let's launch it.
	if enableTLS {
		h := http.Server{
			Addr:        ":443",
			Handler:     s,
			TLSConfig:   tlsConfig,
			ReadTimeout: 30 * time.Second, // use custom timeouts
		}
		if err := h.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			s.Logger.Fatal(err)
		}
	}
	s.Logger.Fatal(s.Start(":" + conf.Port))
}

// defaultJWTSecret sets a random password every time so that
// our endpoints are ALWAYS protected.
func defaultJWTSecret() string {
	return faker.Password()
}
