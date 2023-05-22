package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	httpTimeout = 1 * time.Second

	rootCmd = &cobra.Command{
		Use:   "ff",
		Short: "ff shows 🇨🇦 firearms facts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().Help() //nolint
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
