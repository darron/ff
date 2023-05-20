package redis

import (
	"context"

	"github.com/darron/ff/core"
	"github.com/redis/rueidis"
)

type NewsStoryRepository struct {
	client rueidis.Client
}

func NewNewsStoryRepository(conn string) core.NewsStoryService {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{conn},
		DisableCache: true,
	})
	if err != nil {
		panic(err)
	}
	return NewsStoryRepository{client: client}
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	response, err := nsr.client.Do(ctx, nsr.client.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return &ns, err
	}
	ns, err = core.UnmarshalJSONNewsStory(response)
	return &ns, nil
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	return "not-implimented", nil
}
