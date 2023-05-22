package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send data to HTTP Endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		doSend()
	},
}

func doSend() {
	fmt.Println("Do the HTTP send here")
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
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req
}
