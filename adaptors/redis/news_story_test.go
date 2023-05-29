package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/darron/ff/core"
)

func TestNewsStoryStoreAndFind(t *testing.T) {
	s := miniredis.RunT(t)
	rr := NewsStoryRepository{Conn: s.Addr()}
	r := core.FakeNewsStory()
	r.ID = ""
	id, err := rr.Store(&r)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("ID cannot be blank")
	}
	returnRecord, err := rr.Find(id)
	if err != nil {
		t.Error(err)
	}
	if returnRecord.BodyText != r.BodyText {
		t.Error("Those need to match")
	}
}
