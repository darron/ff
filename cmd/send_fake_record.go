package cmd

import (
	"log"

	"github.com/darron/ff/config"
	"github.com/spf13/cobra"
)

func init() {
	sendFakeCmd.AddCommand(sendFakeRecordCmd)
}

var sendFakeRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Send fake data to HTTP Endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			log.Fatal(err)
		}
		err = sendFakeRecord(conf)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func sendFakeRecord(conf *config.App) error {
	// Create an HTTP client.

	// Get the fake Record

	// Send it over.

	return nil
}
