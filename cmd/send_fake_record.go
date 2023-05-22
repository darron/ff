package cmd

import (
	"log"
	"net/http"
	"net/url"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
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
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), "/records")
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
