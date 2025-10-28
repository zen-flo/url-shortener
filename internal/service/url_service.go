package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zen-flo/url-shortener/internal/model"
)

/*
URLService provides methods for creating, retrieving and deleting shortened URLs.
*/
type URLService struct {
	DB *sqlx.DB
}

/*
NewURLService creates a new instance of URLService with the provided database connection.
*/
func NewURLService(db *sqlx.DB) *URLService {
	return &URLService{DB: db}
}

/*
CreateShortURL generates a short code, saves it in the database and returns the shortened URL record.
*/
func (s *URLService) CreateShortURL(original string) (*model.URL, error) {
	if original == "" {
		return nil, errors.New("original URL cannot be empty")
	}

	short := generateShortCode(6)
	url := &model.URL{
		Original:  original,
		Short:     short,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO urls (original, short, created_at) VALUES (?, ?, ?)`
	result, err := s.DB.Exec(query, url.Original, url.Short, url.CreatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
	}
	url.ID = int(id)

	return url, nil
}

/*
GetOriginalURL retrieves the original URL by its short code.
*/
func (s *URLService) GetOriginalURL(short string) (*model.URL, error) {
	var url model.URL
	err := s.DB.Get(&url, "SELECT * FROM urls WHERE short = ?", short)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

/*
DeleteURL removes a shortened URL from the database by its short code.
*/
func (s *URLService) DeleteURL(short string) error {
	_, err := s.DB.Exec("DELETE FROM urls WHERE short = ?", short)
	return err
}

/*
generateShortCode creates a random, URL-safe short code of given length.
*/
func generateShortCode(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating random bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
