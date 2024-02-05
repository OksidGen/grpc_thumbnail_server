package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type ThumbnailRepository interface {
	SaveThumbnail(videoURL string, thumbnail []byte) error
	GetThumbnail(videoURL string) ([]byte, error)
}

type thumbnailRepository struct {
	db *sqlx.DB
}

func NewThumbnailRepository(db *sqlx.DB) ThumbnailRepository {
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS thumbnails (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    video_url TEXT,
    thumbnail BLOB)`); err != nil {
		log.Fatal().Err(err).Msg("Failed to create thumbnails table")
	}
	return &thumbnailRepository{db: db}
}

func (t *thumbnailRepository) SaveThumbnail(videoURL string, thumbnail []byte) error {
	_, err := t.db.Exec(
		`INSERT INTO thumbnails (video_url, thumbnail) VALUES (?, ?)`,
		videoURL,
		thumbnail,
	)
	return err
}

func (t *thumbnailRepository) GetThumbnail(videoURL string) ([]byte, error) {
	var thumbnail []byte
	err := t.db.Get(&thumbnail, `SELECT thumbnail FROM thumbnails WHERE video_url = ?`, videoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("thumbnail not found")
		}
		return nil, err
	}
	return thumbnail, nil
}
