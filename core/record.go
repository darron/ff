package core

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"gopkg.in/guregu/null.v4"
)

type RecordService interface {
	Find(id string) (*Record, error)
	Store(r *Record) (string, error)
	GetAll() ([]*Record, error)
}

type Record struct {
	ID               string      `json:"id" faker:"-" db:"id"`
	Date             string      `json:"date,omitempty" faker:"year" db:"date"`
	Name             string      `json:"name,omitempty" faker:"first_name_male" db:"name"`
	City             string      `json:"city,omitempty" faker:"oneof: Calgary, Montreal, Vancouver, Toronto" db:"city"`
	Province         string      `json:"province,omitempty" faker:"oneof: AB, QC, BC, ON" db:"province"`
	Licensed         null.Bool   `json:"licensed,omitempty" faker:"nullbool" db:"licensed"`
	Victims          int         `json:"victims,omitempty" faker:"boundary_start=2, boundary_end=10" db:"victims"`
	Deaths           int         `json:"deaths,omitempty" faker:"boundary_start=2, boundary_end=10" db:"deaths"`
	Injuries         int         `json:"injuries,omitempty" faker:"boundary_start=2, boundary_end=10" db:"injuries"`
	Suicide          null.Bool   `json:"suicide,omitempty" faker:"nullbool" db:"suicide"`
	DevicesUsed      string      `json:"devices_used,omitempty" faker:"oneof: Gun, Knife, Pipe, Hands, Axe" db:"devices_used"`
	Firearms         null.Bool   `json:"firearms,omitempty" faker:"nullbool" db:"firearms"`
	PossessedLegally null.Bool   `json:"possessed_legally,omitempty" faker:"-" db:"possessed_legally"`
	Warnings         string      `json:"warnings,omitempty" faker:"sentence" db:"warnings"`
	OICImpact        null.Bool   `json:"oic_impact,omitempty" faker:"nullbool" db:"oic_impact"`
	AISummary        string      `json:"ai_summary,omitempty" faker:"paragraph" db:"ai_summary"`
	NewsStories      []NewsStory `json:"news_stories,omitempty" db:"news_stories"`
}

func UnmarshalJSONRecord(j string) (Record, error) {
	var r Record
	err := json.Unmarshal([]byte(j), &r)
	return r, err
}

func FakeRecord() Record {
	CustomFakerData()
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

func CustomFakerData() {
	_ = faker.AddProvider("nullbool", func(v reflect.Value) (interface{}, error) {
		obj := null.Bool{NullBool: sql.NullBool{
			Bool:  randBool(),
			Valid: true,
		}}
		return obj, nil
	})
}

func randBool() bool {
	rand.Seed(time.Now().UnixNano()) //nolint
	return rand.Float32() < 0.5
}
