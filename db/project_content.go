package db

import (
	"time"
)

//+AR
type ProjectContent struct {
	ID        int       `json:"id" db:"pk"`
	ProjectID string    `json:"project_id"`
	ContentID int       `json:"content_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
