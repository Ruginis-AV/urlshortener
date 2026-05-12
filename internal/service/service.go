package service

import (
	"URL/internal/model"
	"URL/internal/storage"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Service struct {
	store *storage.Storage
}

func New(store *storage.Storage) *Service {
	return &Service{store: store}
}

func generateCode() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder

	for i := 0; i < 6; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (s *Service) CreateLink(originalURL string, expiresAt *time.Time) (model.Link, error) {
	code := generateCode()

	link, err := s.store.CreateLink(originalURL, code, expiresAt)
	if err != nil {
		log.Printf("can't create link %v", err)
		return link, err
	}
	return link, nil
}

func (s *Service) GetLink(code string) (model.Link, error) {
	url, err := s.store.GetByCode(code)
		if err != nil {
			log.Printf("this url is no such %v", err)
			return url, err
		}
		
		if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
			return url, errors.New("link expired")
		}
	err = s.store.IncrementClicks(code)
		if err != nil {
			log.Printf("url is empty %v", err)
		}

	return url, nil
}

func(s *Service) GetStats(code string) (model.Link, error) {
	stat, err := s.store.GetStats(code)
		if err != nil {
			log.Printf("stats is empty %v", err)
			return stat, err
		}
	return stat, nil
}

func(s *Service) Delete(code string) error {
	err := s.store.Delete(code)
		if err != nil {
			log.Printf("tis url is no such %v", err)
			return err
		}
	return nil
}