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
		return &r, err
	}
	defer client.Close()
	return rr.find(ctx, id, client)
}

func (rr RecordRepository) find(ctx context.Context, id string, client *sqlx.DB) (*core.Record, error) {
	r := core.Record{}
	recordQuery := `
		SELECT r.id, r.date, r.name, r.city, r.province, r.licensed, r.victims, r.deaths, r.injuries, r.suicide, r.devices_used, r.firearms,
		r.possessed_legally, r.warnings, r.oic_impact, r.ai_summary, n.id, n.record_id, n.url
		FROM records r
		LEFT JOIN news_stories n ON r.id = n.record_id
		WHERE r.id = ?
	`
	err := client.Select(&r, recordQuery, id)
	fmt.Printf("Record: %#v\n", r)
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
