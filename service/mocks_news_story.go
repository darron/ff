package service

import (
	"errors"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
)

type mockNewsStoryRepository struct{}

func (mrr mockNewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.FakeNewsStory()
	ns.ID = id
	return &ns, nil
}

func (mrr mockNewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	if ns.ID == "" {
		ns.ID = uuid.NewString()
	}
	return ns.ID, nil
}

type mockNewsStoryRepositoryError struct{}

func (mrr mockNewsStoryRepositoryError) Find(id string) (*core.NewsStory, error) {
	ns := core.FakeNewsStory()
	ns.ID = id
	return &ns, errors.New("this is an error")
}

func (mrr mockNewsStoryRepositoryError) Store(ns *core.NewsStory) (string, error) {
	if ns.ID == "" {
		ns.ID = uuid.NewString()
	}
	return ns.ID, errors.New("this is an error")
}
