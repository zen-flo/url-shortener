package model

import "time"

// URL represents a shortened URL
// @name URL
//
//	@example {
//	  "id": 1,
//	  "original": "https://example.com",
//	  "short": "abc123",
//	  "createdAt": "2025-10-30T12:00:00Z"
//	}
type URL struct {
	ID        int       `db:"id" json:"id"`                // Unique identifier
	Original  string    `db:"original" json:"original"`    // Original URL
	Short     string    `db:"short" json:"short"`          // Shortened URL
	CreatedAt time.Time `db:"created_at" json:"createdAt"` // Timestamp when URL was created
}
