package cmd

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
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

func getHTTPClient() *httpclient.Client {
	// Create an HTTP client.
	initalTimeout := 2 * time.Millisecond         // Inital timeout
	maxTimeout := 9 * time.Millisecond            // Max time out
	exponentFactor := 2                           // Multiplier
	maximumJitterInterval := 2 * time.Millisecond // Max jitter interval. It must be more than 1*time.Millisecond

	backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, float64(exponentFactor), maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(httpTimeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(4),
	)
	return client
}

func getHTTPRequest(method, url, body string) *http.Request {
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req
}

func sendFakeRecord(conf *config.App) error {
	client := getHTTPClient()

	u, err := url.JoinPath(conf.GetHTTPEndpoint(), "/records")
	if err != nil {
		return err
	}
	conf.Logger.Debug("sendFakeRecord", "url", u)
	record := core.FakeRecordJSON()
	conf.Logger.Debug("sendFakeRecord", "record", record)
	req := getHTTPRequest(http.MethodPost, u, record)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
