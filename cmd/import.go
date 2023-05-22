package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import CSV with data - exported from Google Sheet",
	Run: func(cmd *cobra.Command, args []string) {
		doImport()
	},
}

func doImport() {
	fmt.Println("Do the import here")
}
