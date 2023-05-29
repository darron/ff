package redis

import (
	"context"
	"encoding/json"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

type RecordRepository struct {
	Conn string
}

func (rr RecordRepository) Connect(conn string) (rueidis.Client, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:       []string{conn},
		DisableCache:      true,
		RingScaleEachConn: 8,
	})
	return client, err
}

func (rr RecordRepository) Find(id string) (*core.Record, error) {
	r := core.Record{}
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	client, err := rr.Connect(rr.Conn)
	if err != nil {
		return &r, err
	}
	defer client.Close()
	response, err := client.Do(ctx, client.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return &r, err
	}
	r, err = core.UnmarshalJSONRecord(response)
	if err != nil {
		return &r, err
	}
	// Get all related NewsStories as well.
	storiesKey := allStoriesPrefix + "-" + id
	storiesSlice, err := client.Do(ctx, client.B().Lrange().Key(storiesKey).Start(0).Stop(-1).Build()).AsStrSlice()
	stories := []core.NewsStory{}
	// Get the NewsStoryRepository - don't really like this.
	nsr := NewsStoryRepository{Conn: rr.Conn} //nolint
	for _, storyID := range storiesSlice {
		story, err := nsr.Find(storyID)
		if err != nil {
			return &r, err
		}
		stories = append(stories, *story)
	}
	r.NewsStories = stories
	return &r, err
}

func (rr RecordRepository) Store(r *core.Record) (string, error) {
	redisKey := "record-" + uuid.NewString()
	r.ID = redisKey
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	client, err := rr.Connect(rr.Conn)
	if err != nil {
		return "", err
	}
	defer client.Close()
	// Remove NewsStories to store them separately.
	newsStories := r.NewsStories
	r.NewsStories = nil
	j, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	err = client.Do(ctx, client.B().Set().Key(redisKey).Value(string(j)).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// Add ID to list of all Records
	err = client.Do(ctx, client.B().Lpush().Key(allRecords).Element(redisKey).Build()).Error()
	if err != nil {
		return redisKey, err
	}
	// Add All NewsStories.
	if len(newsStories) > 0 {
		// Get the NewsStoryRepository - don't really like this.
		nsr := NewsStoryRepository{Conn: rr.Conn} //nolint
		for _, story := range newsStories {
			story.RecordID = redisKey
			_, err := nsr.Store(&story)
			if err != nil {
				return "", err
			}
		}
	}
	return redisKey, err
}

func (rr RecordRepository) GetAll() ([]*core.Record, error) {
	var records []*core.Record

	// Grab the list of IDs stored in allRecords
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	client, err := rr.Connect(rr.Conn)
	if err != nil {
		return records, err
	}
	defer client.Close()
	allRecordIDs, err := client.Do(ctx, client.B().Lrange().Key(allRecords).Start(0).Stop(-1).Build()).AsStrSlice()
	if err != nil {
		return records, err
	}
	// Grab all the Records
	for _, recordID := range allRecordIDs {
		record, err := rr.Find(recordID)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, nil
}
