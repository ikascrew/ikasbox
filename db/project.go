package db

import (
	"time"

	"golang.org/x/xerrors"
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

func NewProject() *Project {
	p := Project{}
	return &p
}

func SelectProject() ([]*Project, error) {
	projects, err := Project{}.All().Query()
	if err != nil {
		return nil, xerrors.Errorf("project all: %w", err)
	}
	return projects, nil

}
