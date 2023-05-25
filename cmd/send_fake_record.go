package cmd

import (
	"log"
	"net/http"
	"net/url"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
)

func init() {
	sendFakeCmd.AddCommand(sendFakeRecordCmd)
}

var sendFakeRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Send fake core.Record to HTTP endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Get(config.WithPort(port), config.WithLogger(logLevel, logFormat))
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
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.RecordsAPIPath)
	if err != nil {
		return err
	}
	conf.Logger.Debug("sendFakeRecord", "url", u)
	record := core.FakeRecordJSON()
	conf.Logger.Debug("sendFakeRecord", "record", record)
	req := getHTTPRequest(http.MethodPost, u, record)
	// Send the HTTP request.
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
