package cmd

import (
	"log"
	"os"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "HTTP Service Commands",
	Run: func(cmd *cobra.Command, args []string) {
		StartService()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

func StartService() {
	// Let's get the config for the app
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	conf.Logger.Info("Starting HTTP Service")
	s, err := service.Get(conf)
	if err != nil {
		conf.Logger.Error(err.Error())
		os.Exit(1)
	}
	s.Logger.Fatal(s.Start(":" + conf.Port))
}
