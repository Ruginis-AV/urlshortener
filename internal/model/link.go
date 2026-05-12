package model

import "time"

type Link struct {
	ID    int
	OriginalURL   string
	Code  string
	CreatedAt time.Time
	ExpiresAt *time.Time
	Clicks int
}