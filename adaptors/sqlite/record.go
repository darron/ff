package sqlite

import (
	"context"
	"fmt"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type RecordRepository struct {
	Filename string
}

func (rr RecordRepository) Connect(filename string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite3", filename)
	return conn, err
}

func (rr RecordRepository) Find(id string) (*core.Record, error) {
	r := core.Record{}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return &r, fmt.Errorf("Find/Connect Error: %w", err)
	}
	defer client.Close()
	return rr.find(ctx, id, client)
}

func (rr RecordRepository) find(ctx context.Context, id string, client *sqlx.DB) (*core.Record, error) {
	r := core.Record{}
	ns := []core.NewsStory{}

	recordQuery := `
		SELECT r.id, r.date, r.name, r.city, r.province, r.licensed, r.victims, r.deaths, r.injuries, r.suicide, r.devices_used, r.firearms,
		r.possessed_legally, r.warnings, r.oic_impact, r.ai_summary
		FROM records r
		WHERE r.id = ?`

	err := client.Get(&r, recordQuery, id)
	if err != nil {
		return &r, fmt.Errorf("find/Get Error: %w", err)
	}
	err = client.Select(&ns, "SELECT * from news_stories WHERE record_id = ?", id)
	if err != nil {
		return &r, fmt.Errorf("find/Select Error: %w", err)
	}
	r.NewsStories = ns
	return &r, err
}

func (rr RecordRepository) Store(record *core.Record) (string, error) {
	id := uuid.NewString()

	// Connect to the db.
	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return "", fmt.Errorf("Store/Connect Error: %w", err)
	}
	defer client.Close()

	// Start the transaction.
	tx, err := client.Begin()
	if err != nil {
		return "", fmt.Errorf("Store/client.Begin Error: %w", err)
	}

	// Insert the Record
	recordQuery := "INSERT INTO records (id, date, name, city, province, licensed, victims, deaths, injuries, suicide, devices_used, firearms, possessed_legally, warnings, oic_impact, ai_summary) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(recordQuery, id, record.Date, record.Name, record.City, record.Province, record.Licensed, record.Victims, record.Deaths, record.Injuries, record.Suicide, record.DevicesUsed, record.Firearms, record.PossessedLegally, record.Warnings, record.OICImpact, record.AISummary)
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("Store/Record/tx.Exec Error: %w", err)
	}

	// Insert the NewsStories
	for _, newsStory := range record.NewsStories {
		storyID := uuid.NewString()
		newsStoryQuery := "INSERT INTO news_stories (id, record_id, url) VALUES (?, ?, ?)"
		_, err = tx.Exec(newsStoryQuery, storyID, id, newsStory.URL)
		if err != nil {
			tx.Rollback() //nolint
			return "", fmt.Errorf("Store/NewsStories/tx.Exec Error: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("Store/tx.Commit Error: %w", err)
	}

	return id, err
}

func (rr RecordRepository) GetAll() ([]*core.Record, error) {
	var ids []string
	var records []*core.Record

	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()

	client, err := rr.Connect(rr.Filename)
	if err != nil {
		return records, fmt.Errorf("GetAll/Connect Error: %w", err)
	}
	defer client.Close()

	// Get all the IDs.
	err = client.Select(&ids, "SELECT id from records")
	if err != nil {
		return records, fmt.Errorf("GetAll/Select/RecordIDs Error: %w", err)
	}

	// Get all of those records with linked stories.
	for _, id := range ids {
		r, err := rr.find(ctx, id, client)
		if err != nil {
			return records, fmt.Errorf("GetAll/find/%s Error: %w", id, err)
		}
		records = append(records, r)
	}

	return records, nil
}
