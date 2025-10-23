package model

import "time"

// URL represents a shortened URL entity.
type URL struct {
	ID        int       `db:"id"`         // Unique identifier
	Original  string    `db:"original"`   // Original URL
	Short     string    `db:"short"`      // Shortened URL
	CreatedAt time.Time `db:"created_at"` // Timestamp when URL was created
}
