package core

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
