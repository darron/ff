package cmd

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/spf13/cobra"
)

var (
	processNewsStoryID          string
	processNewsAll              bool
	processNewsStoryDownloadCmd = &cobra.Command{
		Use:   "news",
		Short: "Download NewsStory body text",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.Get(
				config.WithPort(port),
				config.WithLogger(logLevel, logFormat),
				config.WithJWTToken(jwtToken))
			if err != nil {
				log.Fatal(err)
			}
			if processNewsAll {
				err = processAllNewsStories(conf)
				if err != nil {
					log.Fatal(err)
				} else {
					os.Exit(0)
				}
			}
			err = processNewsStoryDownload(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	processCmd.AddCommand(processNewsStoryDownloadCmd)
	processNewsStoryDownloadCmd.Flags().StringVarP(&processNewsStoryID, "news", "", "", "NewsStort ID to Download NewsStory")
	processNewsStoryDownloadCmd.Flags().BoolVarP(&processNewsAll, "all", "", false, "Download all empty news")
}

func processNewsStoryDownload(conf *config.App) error {
	if processNewsStoryID == "" {
		log.Fatal("Must supply ID for news story")
	}
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.NewsStoriesAPIPathFull, "download", processNewsStoryID)
	if err != nil {
		return err
	}
	conf.Logger.Debug("processNewsStoryDownload", "url", u)
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

func processAllNewsStories(conf *config.App) error {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(10 * time.Minute))
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.NewsStoriesAPIPathFull, "getall")
	if err != nil {
		return err
	}
	conf.Logger.Debug("processAllNewsStories", "url", u)
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
