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

func SelectProject(id int) (*Project, error) {
	project, err := Project{}.Find(id)
	if err != nil {
		return nil, xerrors.Errorf("project Find() error: %w", err)
	}
	return project, nil
}

func SelectProjectList() ([]*Project, error) {
	projects, err := Project{}.All().Query()
	if err != nil {
		return nil, xerrors.Errorf("project All() error: %w", err)
	}
	return projects, nil

}
