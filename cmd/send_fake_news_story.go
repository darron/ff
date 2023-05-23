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
	sendFakeCmd.AddCommand(sendFakeNewsStoryCmd)
}

var sendFakeNewsStoryCmd = &cobra.Command{
	Use:   "news",
	Short: "Send fake core.NewsStory to HTTP endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			log.Fatal(err)
		}
		err = sendFakeNewsStory(conf)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func sendFakeNewsStory(conf *config.App) error {
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), "/stories")
	if err != nil {
		return err
	}
	conf.Logger.Debug("sendFakeNewsStory", "url", u)
	newsStory := core.FakeNewsStoryJSON()
	conf.Logger.Debug("sendFakeNewsStory", "newsStory", newsStory)
	req := getHTTPRequest(http.MethodPost, u, newsStory)
	// Send the HTTP request.
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}