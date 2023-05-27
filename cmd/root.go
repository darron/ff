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
	defaultLogLevel = "debug"
	logLevel        string

	defaultLogFormat = "text"
	logFormat        string

	defaultPort = "8080"
	port        string

	jwtToken string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "level", "", GetENVVariable("LOG_LEVEL", defaultLogLevel), "Log level: info, debug, etc")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "format", "", GetENVVariable("LOG_FORMAT", defaultLogFormat), "Log format: text or json")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", GetENVVariable("PORT", defaultPort), "HTTP Port to listen on")
	rootCmd.PersistentFlags().StringVarP(&jwtToken, "jwt", "", GetENVVariable("JWT_BEARER_TOKEN", ""), "JWT Token to Access API endpoints.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetENVVariable(name, defValue string) string {
	v, ok := os.LookupEnv(name)
	if ok {
		return v
	}
	return defValue
}
