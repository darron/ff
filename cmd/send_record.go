package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/darron/ff/config"
	"github.com/darron/ff/core"
	"github.com/darron/ff/service"
	"github.com/spf13/cobra"
)

func init() {
	sendCmd.AddCommand(sendRecordCmd)
	sendRecordCmd.Flags().StringVarP(&recordDate, "date", "", "", "Record Date: Just numeric year")
	sendRecordCmd.MarkFlagRequired("date") //nolint
	sendRecordCmd.Flags().StringVarP(&recordName, "name", "", "", "Record Name")
	sendRecordCmd.MarkFlagRequired("name") //nolint
	sendRecordCmd.Flags().StringVarP(&recordCity, "city", "", "", "Record City")
	sendRecordCmd.Flags().StringVarP(&recordProvince, "province", "", "", "Record Province")
	sendRecordCmd.MarkFlagRequired("province") //nolint
	sendRecordCmd.Flags().StringVarP(&recordLicensed, "licensed", "", "", "Record licensed"+yesOrNo)
	sendRecordCmd.Flags().IntVarP(&recordVictims, "victims", "", 0, "Record victims - deaths without the perp")
	sendRecordCmd.Flags().IntVarP(&recordDeaths, "deaths", "", 0, "Record deaths - including the perp")
	sendRecordCmd.Flags().IntVarP(&recordInjuries, "injuries", "", 0, "Record injuries")
	sendRecordCmd.Flags().StringVarP(&recordSuicide, "suicide", "", "", "Record Suicide"+yesOrNo)
	sendRecordCmd.Flags().StringVarP(&recordDevices, "devices", "", "", "Record Devices")
	sendRecordCmd.Flags().StringVarP(&recordLegally, "legal", "", "", "Record Legal"+yesOrNo)
	sendRecordCmd.Flags().StringVarP(&recordFirearms, "firearms", "", "", "Record Firearms"+yesOrNo)
	sendRecordCmd.Flags().StringVarP(&recordWarnings, "warnings", "", "", "Record Warnings")
	sendRecordCmd.Flags().StringVarP(&recordOICImpact, "oic", "", "", "Record OIC Impact"+yesOrNo)
}

var (
	yesOrNo         = ": 'Yes' or 'No'"
	recordDate      string
	recordName      string
	recordCity      string
	recordProvince  string
	recordLicensed  string
	recordVictims   int
	recordDeaths    int
	recordInjuries  int
	recordSuicide   string
	recordDevices   string
	recordLegally   string
	recordFirearms  string
	recordWarnings  string
	recordOICImpact string

	sendRecordCmd = &cobra.Command{
		Use:   "record",
		Short: "Send core.Record to HTTP endpoint",
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
			err = sendRecord(conf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func sendRecord(conf *config.App) error {
	client := getHTTPClient()
	// Make up the proper URL including port and path.
	u, err := url.JoinPath(conf.GetHTTPEndpoint(), service.RecordsAPIPathFull)
	if err != nil {
		return err
	}
	conf.Logger.Debug("sendRecord", "url", u)
	r := core.Record{
		Date:             recordDate,
		Name:             recordName,
		City:             recordCity,
		Province:         recordProvince,
		Licensed:         cellToBool(recordLicensed),
		Victims:          recordVictims,
		Deaths:           recordDeaths,
		Injuries:         recordInjuries,
		Suicide:          cellToBool(recordSuicide),
		DevicesUsed:      recordDevices,
		PossessedLegally: cellToBool(recordLegally),
		Firearms:         cellToBool(recordFirearms),
		Warnings:         recordWarnings,
		OICImpact:        cellToBool(recordOICImpact),
	}
	record, _ := json.Marshal(r)
	conf.Logger.Debug("sendRecord", "record", string(record))
	req := getHTTPRequest(http.MethodPost, u, string(record), conf.JWTToken)
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
	return nil
}
