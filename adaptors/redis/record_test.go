package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/darron/ff/core"
)

func TestRecordStoreAndFind(t *testing.T) {
	s := miniredis.RunT(t)
	rr := RecordRepository{Conn: s.Addr()}
	r := core.FakeRecord()
	storiesBefore := len(r.NewsStories)
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
	// Let's make sure the same number of stories are there on both sides.
	storiesAfter := len(returnRecord.NewsStories)
	if storiesBefore != storiesAfter {
		t.Error("We are losing stories on the way in or out")
	}
	if returnRecord.City != r.City {
		t.Error("Those need to match")
	}
}

func TestGetAll(t *testing.T) {
	s := miniredis.RunT(t)
	rr := RecordRepository{Conn: s.Addr()}
	for i := 1; i < 5; i++ {
		r := core.FakeRecord()
		id, err := rr.Store(&r)
		if err != nil {
			t.Error(err)
		}
		if id == "" {
			t.Error("id cannot be blank")
		}
	}
	records, err := rr.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(records) != 4 {
		t.Error("Not enough records in there")
	}
}
