package core

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestRecordFaker(t *testing.T) {
	r := Record{}
	err := faker.FakeData(&r)
	if err != nil {
		t.Error(err)
	}
}

func TestRecordFakerJSON(t *testing.T) {
	j := FakeRecordJSON()
	if j == "" {
		t.Error("JSON was blank")
	}
}

func TestUnmarkshalJSONRecord(t *testing.T) {
	r := FakeRecord()
	j, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}
	r2, err := UnmarshalJSONRecord(string(j))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, r2) {
		t.Error("Those should match")
	}
}
