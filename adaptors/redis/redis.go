package redis

import "time"

var (
	redisTimeout     = time.Second * 4
	allRecords       = "records"
	allStoriesPrefix = "stories"
	recordPrefix     = "record-"
	storyPrefix      = "story-"
)
