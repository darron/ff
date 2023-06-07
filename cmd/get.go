package cmd

import (
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "GET data from HTTP Endpoints",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
}
