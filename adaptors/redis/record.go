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
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{conn}})
	if err != nil {
		panic(err)
	}
	return RecordRepository{client: client}
}

func (rr RecordRepository) Find(id string) (*core.Record, error) {
	r := core.Record{}
	redisKey := "record-" + id
	ctx := context.Background()
	response, err := rr.client.Do(ctx, rr.client.B().Get().Key(redisKey).Build()).ToString()
	if err != nil {
		return &r, err
	}
	r, err = core.UnmarshalJSONRecord(response)
	return &r, err
}

func (rr RecordRepository) Store(ns *core.Record) (string, error) {
	redisKey := "record-" + uuid.NewString()
	ctx := context.Background()
	j, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	err = rr.client.Do(ctx, rr.client.B().Set().Key(redisKey).Value(string(j)).Build()).Error()
	return redisKey, err
}
