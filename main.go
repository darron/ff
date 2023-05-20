package main

import (
	"flag"
	"log"

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
		s, err := service.Get(conf)
		if err != nil {
			log.Fatal(err)
		}
		s.Logger.Fatal(s.Start(":" + conf.Port))
	}
}
