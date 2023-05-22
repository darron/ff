package service

import (
	"errors"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
)

type mockRecordRepository struct{}

func (mrr mockRecordRepository) Find(id string) (*core.Record, error) {
	r := core.FakeRecord()
	r.ID = id
	return &r, nil
}

func (mrr mockRecordRepository) Store(r *core.Record) (string, error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return r.ID, nil
}

type mockRecordRepositoryError struct{}

func (mrr mockRecordRepositoryError) Find(id string) (*core.Record, error) {
	r := core.FakeRecord()
	r.ID = id
	return &r, errors.New("this is an error")
}

func (mrr mockRecordRepositoryError) Store(r *core.Record) (string, error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return r.ID, errors.New("this is an error")
}
