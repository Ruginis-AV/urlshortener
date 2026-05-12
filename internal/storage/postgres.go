package storage

import (
	"database/sql"
	"log"
	"time"

	"URL/internal/model"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) *Storage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic("can't opened database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return &Storage{db: db}
}

func (s *Storage) CreateLink(originalUrl, code string, expiresAt *time.Time) (model.Link, error) {
	link := model.Link{}

	err := s.db.QueryRow(`INSERT INTO links (original_url, code, expires_at) 
						  VALUES ($1, $2, $3) 
						  RETURNING id, created_at`, originalUrl, code, expiresAt).Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		log.Printf("can't insert into links %v", err)
		return link, err
	}

	link.OriginalURL = originalUrl
	link.Code = code
	link.ExpiresAt = expiresAt

	return link, nil
}

func (s *Storage) GetByCode(code string) (model.Link, error) {
	link := model.Link{}
	err := s.db.QueryRow("SELECT * FROM links WHERE code = $1", code).Scan(&link.ID, &link.OriginalURL,
		&link.Code, &link.CreatedAt, &link.ExpiresAt, &link.Clicks)

	if err != nil {
		log.Printf("can't select from links %v", err)
		return link, err
	}
	return link, nil
}

func (s *Storage) IncrementClicks(code string) error {
	_, err := s.db.Exec("UPDATE links SET clicks = clicks+1 WHERE code = $1", code)
	if err != nil {
		log.Printf("can't update links %v", err)
		return err
	}
	return nil
}

func (s *Storage) GetStats(code string) (model.Link, error) {
	link, err := s.GetByCode(code)
	if err != nil {
		log.Printf("can't check stats %v", err)
		return link, err
	}
	return link, nil
}

func (s *Storage) Delete(code string) error {
	_, err := s.db.Exec("DELETE FROM links WHERE code = $1", code)
	if err != nil {
		log.Printf("can't delete string %v", err)
		return err
	}
	return nil
}
