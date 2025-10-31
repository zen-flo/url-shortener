package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zen-flo/url-shortener/internal/model"
)

var (
	urlsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "urls_total",
			Help: "Total number of shortened URLs created.",
		},
	)

	urlsInDB = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "urls_in_db",
			Help: "Current number of shortened URLs stored in the database.",
		},
	)
)

func init() {
	prometheus.MustRegister(urlsTotal)
	prometheus.MustRegister(urlsInDB)
}

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
	s := &URLService{DB: db}
	s.UpdateURLCount()
	return s
}

/*
CreateShortURL generates a unique short code, saves it in the database and returns the shortened URL record.
*/
func (s *URLService) CreateShortURL(original string) (*model.URL, error) {
	if original == "" {
		return nil, errors.New("original URL cannot be empty")
	}

	short := generateShortCode(6)

	// Ensure uniqueness of short code
	for {
		var exists int
		err := s.DB.Get(&exists, "SELECT COUNT(*) FROM urls WHERE short = ?", short)
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			break
		}
		short = generateShortCode(6)
	}

	url := &model.URL{
		Original:  original,
		Short:     short,
		CreatedAt: time.Now(),
	}

	// Insert into database
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

	// Increase Prometheus counter and update gauge
	urlsTotal.Inc()
	s.UpdateURLCount()

	return url, nil
}

/*
GetOriginalURL retrieves the original URL by its short code.
Returns an error if the URL does not exist.
*/
func (s *URLService) GetOriginalURL(short string) (*model.URL, error) {
	var url model.URL
	err := s.DB.Get(&url, "SELECT * FROM urls WHERE short = ?", short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("URL not found")
		}
		return nil, err
	}
	return &url, nil
}

/*
DeleteURL removes a shortened URL from the database by its short code.
Returns an error if the URL does not exist.
*/
func (s *URLService) DeleteURL(short string) error {
	result, err := s.DB.Exec("DELETE FROM urls WHERE short = ?", short)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("URL not found")
	}

	// Update the gauge after deletion
	s.UpdateURLCount()

	return nil
}

/*
UpdateURLCount updates the Prometheus gauge with the current number of URLs in the database.
*/
func (s *URLService) UpdateURLCount() {
	var count int
	err := s.DB.Get(&count, "SELECT COUNT(*) FROM urls")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to count URLs: %v", err)
		}
		return
	}
	urlsInDB.Set(float64(count))
}

/*
generateShortCode creates a random, URL-safe short code of given length.
*/
func generateShortCode(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating random bytes: %v", err)
	}

	// Encode to URL-safe base64 and trim padding
	return base64.URLEncoding.EncodeToString(b)[:length]
}
