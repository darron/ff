package cmd

import (
	"log"
	"os"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/go-faker/faker/v4"
	"github.com/spf13/cobra"
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
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&redisConn, "redisConn", "r", GetENVVariable("REDIS", defaultRedisConn), "Redis connection string")
	serviceCmd.Flags().StringVarP(&jwtSecret, "jwtSecret", "", GetENVVariable("JWT_SECRET", defaultJWTSecret()), "JWT Secret")
}

func StartService() {
	// Setup some options
	var opts []config.OptFunc
	opts = append(opts, config.WithPort(port))
	opts = append(opts, config.WithLogger(logLevel, logFormat))
	// Once we have another db option - we'll adjust this.
	opts = append(opts, config.WithRedis(redisConn))

	// Let's enable JWT if it's defined.
	if jwtSecret != "" {
		opts = append(opts, config.WithJWTSecret(jwtSecret))
	}

	// Let's get the config for the app
	conf, err := config.Get(opts...)
	if err != nil {
		log.Fatal(err)
	}

	conf.Logger.Info("Starting HTTP Service")
	s, err := service.Get(conf)
	if err != nil {
		conf.Logger.Error(err.Error())
		os.Exit(1)
	}

	// Let's get something ready to import data on startup
	// If it's available.
	h := service.HTTPService{}
	go h.ImportOnStartup(conf) //nolint

	s.Logger.Fatal(s.Start(":" + conf.Port))
}

// defaultJWTSecret sets a random password every time so that
// our endpoints are ALWAYS protected.
func defaultJWTSecret() string {
	return faker.Password()
}
