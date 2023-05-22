package cmd

import (
	"log"

	"github.com/darron/ff/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ff",
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func version() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	v := config.GetVersionInfo()
	bsMap := config.GetBuildSettings(v)
	conf.Logger.Info("Version", "go", v.GoVersion, "sha", bsMap["revision"], "time", bsMap["time"], "dirtyTree", bsMap["modified"])
}
