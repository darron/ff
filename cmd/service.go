package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/go-faker/faker/v4"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

	defaultRedisConn = "127.0.0.1:6379"
	redisConn        string

	jwtSecret string

	defaultProfilingEnabled = false
	profilingEnabled        bool

	defaultMiddlewareHTMLCacheEnabled = true
	middlewareHTMLCacheEnabled        bool
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&redisConn, "redisConn", "r", GetENVVariable("REDIS", defaultRedisConn), "Redis connection string")
	serviceCmd.Flags().StringVarP(&jwtSecret, "jwtSecret", "", GetENVVariable("JWT_SECRET", defaultJWTSecret()), "JWT Secret")
	serviceCmd.Flags().BoolVarP(&profilingEnabled, "profiling", "", GetBoolENVVariable("PROFILING_ENABLED", defaultProfilingEnabled), "Enable Datadog tracing and profiling")
	serviceCmd.Flags().BoolVarP(&middlewareHTMLCacheEnabled, "htmlcache", "", GetBoolENVVariable("HTMLCACHE_ENABLED", defaultMiddlewareHTMLCacheEnabled), "Enable Middleware Cache")
}

func StartService() {
	// Setup some options
	var opts []config.OptFunc
	opts = append(opts, config.WithPort(port))
	opts = append(opts, config.WithLogger(logLevel, logFormat))
	opts = append(opts, config.WithMiddlewareHTMLCache(middlewareHTMLCacheEnabled))

	// Once we have another db option - we'll adjust this.
	opts = append(opts, config.WithRedis(redisConn))

	// Let's enable JWT if it's defined.
	if jwtSecret != "" {
		opts = append(opts, config.WithJWTSecret(jwtSecret))
	}

	// Let's turn on TLS.
	if enableTLS {
		tls := config.TLS{
			CacheDir:    tlsCache,
			DomainNames: tlsDomains,
			Email:       tlsEmail,
			Enable:      enableTLS,
		}
		err := tls.Verify()
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, config.WithTLS(tls))
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

	// If we are going to turn on TLS - let's launch it.
	if enableTLS {
		domains := strings.Split(conf.TLS.DomainNames, ",")
		autoTLSManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(conf.TLS.CacheDir),
			Email:      conf.TLS.Email,
			HostPolicy: autocert.HostWhitelist(domains...),
		}
		h := http.Server{
			Addr:    ":443",
			Handler: s,
			TLSConfig: &tls.Config{
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			},
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
