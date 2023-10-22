package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
)

func init() {
	sendCmd.AddCommand(sendNewsStoryCmd)
	sendNewsStoryCmd.Flags().StringVarP(&recordID, "record", "", "", "Record ID to add News Story To")
	sendNewsStoryCmd.Flags().StringVarP(&newsURL, "news", "", "", "News story URL to add to Record")
}

var (
	sendNewsStoryCmd = &cobra.Command{
		Use:   "news",
		Short: "Send core.NewsStory to HTTP endpoint",
		Run: func(cmd *cobra.Command, args []string) {
			tlsConfig := config.TLS{
				DomainNames: tlsDomains,
			}
			conf, err := config.Get(
				config.WithPort(port),
				config.WithLogger(logLevel, logFormat),
				config.WithJWTToken(jwtToken),
				config.WithTLS(tlsConfig),
			)
			if err != nil {
				log.Fatal(err)
			}
			err = sendNewsStory(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	recordID string
	newsURL  string
)

func sendNewsStory(conf *config.App) error {
	if recordID == "" {
		log.Fatal("Need Record id to send add News to")
	}
	if newsURL == "" {
		log.Fatal("Need News URL to add to Record")
	}
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.NewsStoriesAPIPathFull)
	if err != nil {
		return err
	}
	conf.Logger.Debug("sendNewsStory", "url", u)

	ns := core.NewsStory{
		RecordID: recordID,
		URL:      newsURL,
	}
	newsStory, _ := json.Marshal(ns)
	conf.Logger.Debug("sendNewsStory", "newsStory", string(newsStory))
	req := getHTTPRequest(http.MethodPost, u, string(newsStory), conf.JWTToken)
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
