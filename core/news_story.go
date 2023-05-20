package core

import "encoding/json"

type NewsStoryService interface {
	Find(id string) (*NewsStory, error)
	Store(ns *NewsStory) (string, error)
}

type NewsStory struct {
	ID        string
	RecordID  string
	URL       string
	BodyText  string
	AISummary string
}

func UnmarshalJSONNewsStory(j string) (NewsStory, error) {
	var ns NewsStory
	err := json.Unmarshal([]byte(j), &ns)
	return ns, err
}
