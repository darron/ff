package cmd

import (
	"fmt"

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
