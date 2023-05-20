package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/darron/ff/core"
)

func TestStoreAndFind(t *testing.T) {
	s := miniredis.RunT(t)
	rr := NewRecordRepository(s.Addr())
	r := core.FakeRecord()
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
	if returnRecord.City != r.City {
		t.Error("Those need to match")
	}
}
