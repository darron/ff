package cmd

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/darron/ff/config"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getNewsStoryCmd)
	getNewsStoryCmd.Flags().StringVarP(&getNewsStoryID, "news", "", "", "News Story ID to GET")
}

var (
	getNewsStoryID  string
	getNewsStoryCmd = &cobra.Command{
		Use:   "news",
		Short: "Get core.NewsStory from HTTP endpoint",
		Run: func(cmd *cobra.Command, args []string) {
			conf, err := config.Get(
				config.WithPort(port),
				config.WithLogger(logLevel, logFormat),
				config.WithJWTToken(jwtToken))
			if err != nil {
				log.Fatal(err)
			}
			err = getNewsStory(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func getNewsStory(conf *config.App) error {
	if getNewsStoryID == "" {
		conf.Logger.Error("Must pass along news story ID")
		os.Exit(1)
	}
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.NewsStoriesAPIPathFull, getNewsStoryID)
	if err != nil {
		return err
	}
	conf.Logger.Debug("getNewsStory", "url", u)
	req := getHTTPRequest(http.MethodGet, u, "", conf.JWTToken)
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
