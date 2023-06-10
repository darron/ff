package sqlite

import (
	"context"
	"database/sql"

	"github.com/darron/ff/core"
	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

type RecordRepository struct {
	Filename string
}

func (rr RecordRepository) Connect(filename string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", filename)
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

func (rr RecordRepository) Store(record *core.Record) (string, error) {
	id := uuid.NewString()

	// Connect to the db.
	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Start the transaction.
	tx, err := client.Begin()
	if err != nil {
		return "", err
	}

	// Insert the Record
	recordQuery := "INSERT INTO records (id, date, name, city, province, licensed, victims, deaths, injuries, suicide, devices_used, firearms, possessed_legally, warnings, oic_impact, ai_summary) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(recordQuery, id, record.Date, record.Name, record.City, record.Province, record.Licensed, record.Victims, record.Deaths, record.Injuries, record.Suicide, record.DevicesUsed, record.Firearms, record.PossessedLegally, record.Warnings, record.OICImpact, record.AISummary)
	if err != nil {
		tx.Rollback() //nolint
		return "", err
	}

	// Insert the NewsStories
	for _, newsStory := range record.NewsStories {
		newsStoryQuery := "INSERT INTO news_stories (record_id, url) VALUES (?, ?)"
		_, err = tx.Exec(newsStoryQuery, id, newsStory.URL)
		if err != nil {
			tx.Rollback() //nolint
			return "", err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback() //nolint
		return "", err
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