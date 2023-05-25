package core

import (
	"encoding/json"

	"github.com/go-faker/faker/v4"
	"gopkg.in/guregu/null.v4"
)

type RecordService interface {
	Find(id string) (*Record, error)
	Store(r *Record) (string, error)
	GetAll() ([]*Record, error)
}

type Record struct {
	ID               string      `json:"id" faker:"-"`
	Date             string      `json:"date,omitempty" faker:"date"`
	Name             string      `json:"name,omitempty" faker:"first_name_male"`
	City             string      `json:"city,omitempty" faker:"oneof: Calgary, Montreal, Vancouver, Toronto"`
	Province         string      `json:"province,omitempty" faker:"oneof: AB, QC, BC, ON"`
	Licensed         null.Bool   `json:"licensed,omitempty" faker:"-"`
	Victims          int         `json:"victims,omitempty" faker:"boundary_start=2, boundary_end=10"`
	Deaths           int         `json:"deaths,omitempty" faker:"boundary_start=2, boundary_end=10"`
	Injuries         int         `json:"injuries,omitempty" faker:"boundary_start=2, boundary_end=10"`
	Suicide          null.Bool   `json:"suicide,omitempty" faker:"-"`
	DevicesUsed      string      `json:"devicesused,omitempty" faker:"oneof: Gun, Knife, Pipe, Hands, Axe"`
	Firearms         null.Bool   `json:"firearms,omitempty" faker:"-"`
	PossessedLegally null.Bool   `json:"possessedlegally,omitempty" faker:"-"`
	Warnings         string      `json:"warnings,omitempty" faker:"sentence"`
	OICImpact        null.Bool   `json:"oicimpact,omitempty" faker:"-"`
	AISummary        string      `json:"aisummary,omitempty" faker:"paragraph"`
	NewsStories      []NewsStory `json:"news_stories,omitempty"`
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

func FakeRecordJSON() string {
	r := FakeRecord()
	j, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(j)
}
