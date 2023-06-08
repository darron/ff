package sqlite

import (
	"context"
	"database/sql"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
)

type RecordRepository struct {
	Filename string
}

func (rr RecordRepository) Connect(filename string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3_extended", filename)
	return conn, err
}

func (rr RecordRepository) Find(id string) (*core.Record, error) {
	r := core.Record{}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return &r, err
	}
	defer client.Close()
	return rr.find(ctx, id, client)
}

func (rr RecordRepository) find(ctx context.Context, id string, client *sql.DB) (*core.Record, error) {
	r := core.Record{}
	// TODO: Get all related news stories.
	// TODO: Deal with rows
	_, err := client.Query("SELECT * from records WHERE ID = ?", id)
	if err != nil {
		return &r, err
	}
	return &r, err
}

func (rr RecordRepository) Store(r *core.Record) (string, error) {
	id := uuid.NewString()
	// TODO: When we do the query - we'll use this.
	// ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	// defer cancel()
	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return "", err
	}
	defer client.Close()
	// Remove NewsStories to store them separately.
	newsStories := r.NewsStories
	r.NewsStories = nil

	// TODO: Insert record
	_, err = client.Query("INSERT into records = ?", id)

	// TODO: Insert all news stories.
	for _, story := range newsStories {
		_, err = client.Query("INSERT into records = ?", story)
	}

	return id, err
}

func (rr RecordRepository) GetAll() ([]*core.Record, error) {
	var records []*core.Record

	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return records, err
	}
	defer client.Close()

	// TODO: Do a select with a join as well.
	// TODO: Then deal with all the rows.
	_, err = client.Query("SELECT * from records")
	if err != nil {
		return records, err
	}

	return records, nil
}
