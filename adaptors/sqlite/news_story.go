package sqlite

import (
	"context"
	"fmt"

	"github.com/darron/ff/core"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/glebarez/go-sqlite"
)

type NewsStoryRepository struct {
	Filename string
}

func (nsr NewsStoryRepository) Connect(filename string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("sqlite3", filename)
	return conn, err
}

func (nsr NewsStoryRepository) Find(id string) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return &ns, fmt.Errorf("Find/Connect Error: %w", err)
	}
	defer client.Close()
	return nsr.find(ctx, id, client)
}

func (nsr NewsStoryRepository) find(ctx context.Context, id string, client *sqlx.DB) (*core.NewsStory, error) {
	ns := core.NewsStory{}
	err := client.Get(&ns, "SELECT * from news_stories WHERE id = ?", id)
	if err != nil {
		return &ns, fmt.Errorf("find/Get Error: %w", err)
	}
	return &ns, err
}

func (nsr NewsStoryRepository) Store(ns *core.NewsStory) (string, error) {
	if ns.ID == "" {
		ns.ID = uuid.NewString()
	}
	ctx, cancel := context.WithTimeout(context.Background(), sqliteTimeout)
	defer cancel()
	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return "", fmt.Errorf("Store/Connect Error: %w", err)
	}
	defer client.Close()

	return nsr.store(ctx, ns, client)
}

func (nsr NewsStoryRepository) store(ctx context.Context, ns *core.NewsStory, client *sqlx.DB) (string, error) {
	// Start the transaction.
	tx, err := client.Begin()
	if err != nil {
		return "", fmt.Errorf("store/client.Begin Error: %w", err)
	}

	// Insert/Upsert the NewsStory
	newsStoryQuery := `INSERT INTO news_stories (id, record_id, url, body_text, ai_summary) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id)
		DO UPDATE SET url=excluded.url, body_text=excluded.body_text, ai_summary=excluded.ai_summary`
	_, err = tx.Exec(newsStoryQuery, ns.ID, ns.RecordID, ns.URL, ns.BodyText, ns.AISummary)
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("store/NewsStories/tx.Exec Error: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback() //nolint
		return "", fmt.Errorf("store/tx.Commit Error: %w", err)
	}

	return ns.ID, err
}

func (nsr NewsStoryRepository) GetAll() ([]*core.NewsStory, error) {
	var stories []*core.NewsStory

	client, err := nsr.Connect(nsr.Filename)
	if err != nil {
		return stories, fmt.Errorf("GetAll/Connect Error: %w", err)
	}
	defer client.Close()

	err = client.Select(&stories, "SELECT * from news_stories")
	if err != nil {
		return stories, fmt.Errorf("GetAll/Select Error: %w", err)
	}

	return stories, nil
}
