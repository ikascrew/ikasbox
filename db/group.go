package db

import (
	"time"
)

const CreateGroupsSQL = `
CREATE TABLE [GROUPS] (
    [id] INTEGER PRIMARY KEY AUTOINCREMENT,
    [name] VARCHAR(64) NOT NULL,
    [path] VARCHAR(512) DEFAULT '',
    [created_at] DATETIME,
    [updated_at] DATETIME
)`

//+AR
type Group struct {
	ID        int       `json:"id" db:"pk"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SelectGroup() ([]*Group, error) {
	groups, err := Group{}.Order("id", "asc").All().Query()
	if err != nil {
		return nil, err
	}
	return groups, err
}

func FindGroup(id int) (*Group, error) {
	return Group{}.Find(id)
}
