package cmd

import (
	"github.com/spf13/cobra"
)

var (
	processCmd = &cobra.Command{
		Use:   "process",
		Short: "Process things as necessary",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(processCmd)
}
