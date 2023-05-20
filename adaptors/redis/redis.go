package redis

import "time"

var (
	redisTimeout     = time.Second * 2
	allRecords       = "records"
	allStoriesPrefix = "stories"
)
