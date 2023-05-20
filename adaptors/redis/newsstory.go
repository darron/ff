package redis

import (
	"github.com/darron/ff/core"
	"github.com/redis/rueidis"
)

type NewsStoryRepository struct {
	client rueidis.Client
}

func NewNewsStoryRepository(conn string) core.NewsStoryService {
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{conn}})
	if err != nil {
		panic(err)
	}
	return NewsStoryRepository{client: client}
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	return &ns, nil
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	return "not-implimented", nil
}
