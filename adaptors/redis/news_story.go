package redis

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

type NewsStoryRepository struct {
	Conn string
}

func (nsr NewsStoryRepository) Connect(conn string) (rueidis.Client, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:       []string{conn},
		DisableCache:      true,
		RingScaleEachConn: 8,
	})
	return client, err
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Conn)
	if err != nil {
		return &ns, err
	}
	defer client.Close()
	return nsr.find(ctx, id, client)
}

func (nsr NewsStoryRepository) find(ctx context.Context, id string, client rueidis.Client) (*core.NewsStory, error) {
	if !strings.HasPrefix(id, storyPrefix) {
		id = storyPrefix + id
	}
	ns := core.NewsStory{}
	response, err := client.Do(ctx, client.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return &ns, err
	}
	ns, err = core.UnmarshalJSONNewsStory(response)
	return &ns, err
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	redisKey := storyPrefix + uuid.NewString()
	ns.ID = redisKey
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	j, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	client, err := nsr.Connect(nsr.Conn)
	if err != nil {
		return "", err
	}
	defer client.Close()
	err = client.Do(ctx, client.B().Set().Key(redisKey).Value(string(j)).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// Add ID to list of all Stories for a record.
	storyList := allStoriesPrefix + "-" + ns.RecordID
	err = client.Do(ctx, client.B().Lpush().Key(storyList).Element(redisKey).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	return redisKey, err
}
