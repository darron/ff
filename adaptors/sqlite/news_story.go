package sqlite

import (
	"context"
	"database/sql"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
)

type NewsStoryRepository struct {
	Filename string
}

func (nsr NewsStoryRepository) Connect(filename string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3_extended", filename)
	return conn, err
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	r := core.NewsStory{}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return &r, err
	}
	defer client.Close()
	return nsr.find(ctx, id, client)
}

func (nsr NewsStoryRepository) find(ctx context.Context, id string, client *sql.DB) (*core.NewsStory, error) {
	r := core.NewsStory{}
	// TODO: Get all related news stories.
	// TODO: Deal with rows
	_, err := client.Query("SELECT * from news_stories WHERE record_id = ?", id)
	if err != nil {
		return &r, err
	}
	return &r, err
}

func (nsr NewsStoryRepository) Store(r *core.NewsStory) (string, error) {
	id := uuid.NewString()
	// TODO: When we do the query - we'll use this.
	// ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	// defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// TODO:
	_, err = client.Query("INSERT into news_story = ?", r)

	return id, err
}

func (nsr NewsStoryRepository) GetAll() ([]*core.NewsStory, error) {
	var stories []*core.NewsStory

	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return stories, err
	}
	defer client.Close()

	// TODO: Do a select with a join as well.
	// TODO: Then deal with all the rows.
	_, err = client.Query("SELECT * from news_stories")
	if err != nil {
		return stories, err
	}

	return stories, nil
}
