package cmd

import (
	"github.com/spf13/cobra"
)

var (
	summarizeCmd = &cobra.Command{
		Use:   "summarize",
		Short: "Summarize things using OpenAI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
}
