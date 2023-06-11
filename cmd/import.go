package cmd

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
	"gopkg.in/guregu/null.v4"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&importFilename, "import", "i", GetENVVariable("IMPORT_FILENAME", defaultImportFilename), "Filename to import")
}

var (
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import CSV with data - exported from Google Sheet",
		Run: func(cmd *cobra.Command, args []string) {
			var opts []config.OptFunc
			opts = append(opts, config.WithPort(port))
			opts = append(opts, config.WithLogger(logLevel, logFormat))
			opts = append(opts, config.WithJWTToken(jwtToken))
			// We only need the right domain name to connect with
			// TLS - we don't need any of the other values.
			if enableTLS && tlsDomains != "" {
				tls := config.TLS{
					DomainNames: tlsDomains,
					Enable:      enableTLS,
				}
				opts = append(opts, config.WithTLS(tls))
			}
			conf, err := config.Get(opts...)
			if err != nil {
				log.Fatal(err)
			}
			err = doImport(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	importFilename        string
	defaultImportFilename = "import.csv"
)

func doImport(conf *config.App) error {
	// Read in the CSV.
	data, err := os.ReadFile(importFilename)
	if err != nil {
		return err
	}
	r := csv.NewReader(bytes.NewReader(data))
	lines, err := r.ReadAll()
	if err != nil {
		return err
	}
	// Header:
	// Date,Name,Location,Province,Licensed,Victims,Deaths,Injuries,Suicide,Devices Used,Posessed Legally?,Firearms,Warnings,OIC Impact,Links,,,,,,,,,
	// For each line:
	// 	Convert each line to core.Record
	// 	After OIC Impact - each additional non blank entry is a core.NewsStory.
	// 	Send the core.Record to service.RecordsAPIPath
	for _, line := range lines {
		r := core.Record{
			Date:             line[0],
			Name:             line[1],
			City:             (strings.Split(line[2], ","))[0],
			Province:         line[3],
			Licensed:         cellToBool(line[4]),
			Victims:          cellToInt(line[5]),
			Deaths:           cellToInt(line[6]),
			Injuries:         cellToInt(line[7]),
			Suicide:          cellToBool(line[8]),
			DevicesUsed:      line[9],
			PossessedLegally: cellToBool(line[10]),
			Firearms:         cellToBool(line[11]),
			Warnings:         line[12],
			OICImpact:        cellToBool(line[13]),
		}
		// Let's deal with the news stories
		var stories []core.NewsStory
		links := line[14:]
		for _, link := range links {
			if link != "" {
				ns := core.NewsStory{}
				ns.URL = cleanUpLink(link)
				stories = append(stories, ns)
			}
		}
		r.NewsStories = stories

		// Setup the HTTP client:
		client := getHTTPClient()
		// Make up the proper URL including port and path.
		u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.RecordsAPIPathFull)
		if err != nil {
			return err
		}
		jsonRecord, err := json.Marshal(r)
		if err != nil {
			return err
		}
		req := getHTTPRequest(http.MethodPost, u, string(jsonRecord), conf.JWTToken)
		// Send the HTTP request.
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		conf.Logger.Info(string(body))
	}
	return nil
}

func cleanUpLink(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	q := u.Query()
	if q.Has("fbclid") {
		q.Del("fbclid")
		u.RawQuery = q.Encode()
		return u.String()
	}
	return link
}

func cellToInt(cell string) int {
	if cell == "" {
		return 0
	}
	i, err := strconv.Atoi(cell)
	if err != nil {
		return 0
	}
	return i
}

func cellToBool(cell string) null.Bool {
	switch cell {
	case "Yes":
		return null.Bool{NullBool: sql.NullBool{
			Bool:  true,
			Valid: true,
		}}
	case "No":
		return null.Bool{NullBool: sql.NullBool{
			Bool:  false,
			Valid: true,
		}}
	default:
		return null.Bool{NullBool: sql.NullBool{
			Bool:  false,
			Valid: false,
		}}
	}
}
