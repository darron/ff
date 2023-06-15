package cmd

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/spf13/cobra"
)

var (
	summarizeNewsStoryCmd = &cobra.Command{
		Use:   "news",
		Short: "Summarize News Story using OpenAI",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.Get(
				config.WithPort(port),
				config.WithLogger(logLevel, logFormat),
				config.WithJWTToken(jwtToken))
			if err != nil {
				log.Fatal(err)
			}
			err = processNewsStorySummary(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	summarizeNewsStoryID string
)

func init() {
	summarizeCmd.AddCommand(summarizeNewsStoryCmd)
	summarizeNewsStoryCmd.Flags().StringVarP(&summarizeNewsStoryID, "news", "", "", "NewsStory ID to Summarize")
}

func processNewsStorySummary(conf *config.App) error {
	if summarizeNewsStoryID == "" {
		log.Fatal("Must supply ID for news story")
	}
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(10 * time.Minute))
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.NewsStoriesAPIPathFull, "summarize", summarizeNewsStoryID)
	if err != nil {
		return err
	}
	conf.Logger.Debug("processNewsStorySummary", "url", u)
	req := getHTTPRequest(http.MethodPost, u, "", conf.JWTToken)
	// Send the HTTP request.
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	conf.Logger.Info(string(body))
	return nil
}
