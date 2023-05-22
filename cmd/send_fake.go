package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendFakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Send fake data to HTTP Endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pick an option - record or news")
	},
}

func init() {
	sendCmd.AddCommand(sendFakeCmd)
}
