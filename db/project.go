package db

import (
	"time"
)

const CreateProjectsSQL = `
CREATE TABLE [PROJECTS] (
    [id] INTEGER PRIMARY KEY AUTOINCREMENT,
    [name] TEXT,
    [width] INTEGER,
    [height] INTEGER,
    [default_content] INTEGER,
    [created_at] DATETIME,
    [updated_at] DATETIME
)
`

//+AR
type Project struct {
	ID             int       `json:"id" db:"pk"`
	Name           string    `json:"name"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	DefaultContent int       `json:"default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func SelectProject() ([]*Project, error) {
	projects, err := Project{}.All().Query()
	if err != nil {
		return nil, err
	}
	return projects, nil

}
