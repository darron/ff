package core

type RecordService interface {
	Find(id string) (*Record, error)
	Store(r *Record) (string, error)
}

type NewsStoryService interface {
	Find(id string) (*NewsStory, error)
	Store(ns *NewsStory) (string, error)
}
