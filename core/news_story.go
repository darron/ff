package core

import (
	"encoding/json"

	"github.com/go-faker/faker/v4"
)

type NewsStoryService interface {
	Find(id string) (*NewsStory, error)
	Store(ns *NewsStory) (string, error)
}

type NewsStory struct {
	ID        string `json:"id" faker:"-" sql:"id"`
	RecordID  string `json:"record_id,omitempty" faker:"uuid_hyphenated" sql:"record_id"`
	URL       string `json:"url,omitempty" faker:"url" sql:"url"`
	BodyText  string `json:"body_text,omitempty" faker:"paragraph" sql:"body_text"`
	AISummary string `json:"ai_summary,omitempty" faker:"sentence" sql:"ai_summary"`
}

func UnmarshalJSONNewsStory(j string) (NewsStory, error) {
	var ns NewsStory
	err := json.Unmarshal([]byte(j), &ns)
	return ns, err
}

func FakeNewsStory() NewsStory {
	ns := NewsStory{}
	faker.FakeData(&ns) //nolint
	return ns
}

func FakeNewsStoryJSON() string {
	ns := FakeNewsStory()
	j, err := json.Marshal(ns)
	if err != nil {
		return ""
	}
	return string(j)
}
