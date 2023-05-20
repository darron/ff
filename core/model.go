package core

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID               string    `json:"id"`
	Date             time.Time `json:"date,omitempty"`
	Name             string    `json:"name,omitempty"`
	City             string    `json:"city,omitempty"`
	Province         string    `json:"province,omitempty"`
	Licensed         bool      `json:"licensed,omitempty"`
	Victims          int       `json:"victims,omitempty"`
	Deaths           int       `json:"deaths,omitempty"`
	Injuries         int       `json:"injuries,omitempty"`
	Suicide          bool      `json:"suicide,omitempty"`
	DevicesUsed      string    `json:"devicesused,omitempty"`
	Firearms         bool      `json:"firearms,omitempty"`
	PossessedLegally bool      `json:"possessedlegally,omitempty"`
	Warnings         string    `json:"warnings,omitempty"`
	OICImpact        bool      `json:"oicimpact,omitempty"`
	AISummary        string    `json:"aisummary,omitempty"`
}

func UnmarshalJSONRecord(j string) (Record, error) {
	var r Record
	err := json.Unmarshal([]byte(j), &r)
	return r, err
}

type NewsStory struct {
	ID        string
	RecordID  string
	URL       string
	BodyText  string
	AISummary string
}
