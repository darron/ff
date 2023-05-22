package redis

import (
	"context"
	"encoding/json"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
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
	return &ns, err
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	redisKey := "story-" + uuid.NewString()
	ns.ID = redisKey
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	j, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	err = nsr.client.Do(ctx, nsr.client.B().Set().Key(redisKey).Value(string(j)).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// Add ID to list of all Stories for a record.
	storyList := allStoriesPrefix + "-" + ns.RecordID
	err = nsr.client.Do(ctx, nsr.client.B().Lpush().Key(storyList).Element(redisKey).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	return redisKey, err
}
