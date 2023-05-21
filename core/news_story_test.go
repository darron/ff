package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFakeNewsStory(t *testing.T) {
	j := FakeNewsStoryJSON()
	if j == "" {
		t.Error("JSON was blank")
	}
}

func TestUnmarkshalJSONNewsStory(t *testing.T) {
	r := FakeNewsStory()
	j, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}
	r2, err := UnmarshalJSONNewsStory(string(j))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, r2) {
		t.Error("Those should match")
	}
}
