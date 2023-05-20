package main

import (
	"flag"
	"log"

	"github.com/darron/ff/cmd/service"
)

var (
	defaultStartHTTPService = true
)

func main() {
	var startHTTPService bool
	flag.BoolVar(&startHTTPService, "start", defaultStartHTTPService, "Start HTTP Service")
	flag.Parse()

	if startHTTPService {
		err := service.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}
