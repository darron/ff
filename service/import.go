package service

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/darron/ff/config"
)

var (
	importSleep = 2 * time.Second
	importFile  = "import.csv"
)

func (s HTTPService) ImportOnStartup(conf *config.App) error {
	// Let's wait a little bit.
	time.Sleep(importSleep)

	// Are there any Records in the DB?
	// If there are - let's do nothing and exit.
	records, err := conf.RecordRepository.GetAll()
	if err != nil {
		return err
	}
	if len(records) > 0 {
		return nil
	}

	// Is there a file nearby called importFile?
	// If so - let's import it.
	if _, err := os.Stat(importFile); err == nil {
		conf.Logger.Info("Import file exists - importing")

		// Do the import.
		cmd := exec.Cmd{
			Path: "/",
			Args: []string{"ff", "import"},
			Env: []string{
				"PORT=" + conf.Port,
				"LOG_LEVEL=info",
				"LOG_FORMAT=text",
				"JWT_BEARER_TOKEN=" + conf.JWTToken},
		}
		stdoutStderr, err := cmd.CombinedOutput()
		fmt.Printf("%s\n", stdoutStderr)
		return err
	}

	return nil
}
