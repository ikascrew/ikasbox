package db

import (
	"time"
)

//+AR
type Project struct {
	ID        string    `json:"id" db:"pk"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Default   int       `json:"default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SelectProject() ([]*Project, error) {
	projects, err := Project{}.All().Query()
	if err != nil {
		return nil, err
	}
	return projects, nil

}
