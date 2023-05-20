package redis

import (
	"context"
	"encoding/json"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

type RecordRepository struct {
	client rueidis.Client
}

func NewRecordRepository(conn string) core.RecordService {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{conn},
		DisableCache: true,
	})
	if err != nil {
		panic(err)
	}
	return RecordRepository{client: client}
}

func (rr RecordRepository) Find(id string) (*core.Record, error) {
	r := core.Record{}
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	response, err := rr.client.Do(ctx, rr.client.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return &r, err
	}
	r, err = core.UnmarshalJSONRecord(response)
	// TODO: Get all related NewsStories as well.
	return &r, err
}

func (rr RecordRepository) Store(ns *core.Record) (string, error) {
	redisKey := "record-" + uuid.NewString()
	ns.ID = redisKey
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	j, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	err = rr.client.Do(ctx, rr.client.B().Set().Key(redisKey).Value(string(j)).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// Add ID to list of all Records
	err = rr.client.Do(ctx, rr.client.B().Lpush().Key(allStories).Element(redisKey).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// TODO: Add NewsStories.
	return redisKey, err
}
