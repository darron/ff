package core

import (
	"encoding/json"

	"github.com/go-faker/faker/v4"
)

type RecordService interface {
	Find(id string) (*Record, error)
	Store(r *Record) (string, error)
}

type Record struct {
	ID               string `json:"id" faker:"-"`
	Date             string `json:"date,omitempty" faker:"date"`
	Name             string `json:"name,omitempty" faker:"first_name_male"`
	City             string `json:"city,omitempty" faker:"oneof: Calgary, Montreal, Vancouver, Toronto"`
	Province         string `json:"province,omitempty" faker:"oneof: AB, QC, BC, ON"`
	Licensed         bool   `json:"licensed,omitempty" faker:"-"`
	Victims          int    `json:"victims,omitempty" faker:"oneof: 1, 2"`
	Deaths           int    `json:"deaths,omitempty" faker:"oneof: 1, 2"`
	Injuries         int    `json:"injuries,omitempty" faker:"oneof: 1, 2"`
	Suicide          bool   `json:"suicide,omitempty" faker:"-"`
	DevicesUsed      string `json:"devicesused,omitempty" faker:"oneof: Gun, Knife, Pipe, Hands, Axe"`
	Firearms         bool   `json:"firearms,omitempty" faker:"-"`
	PossessedLegally bool   `json:"possessedlegally,omitempty" faker:"-"`
	Warnings         string `json:"warnings,omitempty" faker:"sentence"`
	OICImpact        bool   `json:"oicimpact,omitempty" faker:"-"`
	AISummary        string `json:"aisummary,omitempty" faker:"paragraph"`
}

func UnmarshalJSONRecord(j string) (Record, error) {
	var r Record
	err := json.Unmarshal([]byte(j), &r)
	return r, err
}

func FakeRecord() Record {
	r := Record{}
	faker.FakeData(&r) //nolint
	return r
}
