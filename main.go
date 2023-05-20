package main

import (
	"flag"
	"log"
	"os"

	"github.com/darron/ff/cmd/service"
	"github.com/darron/ff/config"
)

var (
	defaultStartHTTPService = true
)

func main() {
	var startHTTPService bool
	flag.BoolVar(&startHTTPService, "start", defaultStartHTTPService, "Start HTTP Service")
	flag.Parse()

	// Let's get the config for the app
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if startHTTPService {
		conf.Logger.Info("Starting HTTP Service")
		s, err := service.Get(conf)
		if err != nil {
			conf.Logger.Error(err.Error())
			os.Exit(1)
		}
		s.Logger.Fatal(s.Start(":" + conf.Port))
	}
}
