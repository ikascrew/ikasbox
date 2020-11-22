package db

import (
	"time"
)

const (
	ContentPageNum    = 100
	CreateContentsSQL = `
CREATE TABLE [CONTENTS] (
    [id] INTEGER PRIMARY KEY AUTOINCREMENT,
    [group_id] INTEGER,
    [name] VARCHAR(128) NOT NULL,
    [type] VARCHAR(32),
    [path] VARCHAR(1024),
    [width] INTEGER,
    [height] INTEGER,
    [fps] REAL,
    [frames] INTEGER,
    [fourcc] REAL,
    [created_at] DATETIME,
    [updated_at] DATETIME
)
`
)

//+AR
type Content struct {
	ID        int       `json:"id" db:"pk"`
	GroupId   int       `json:"group_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Path      string    `json:"path"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	FPS       float64   `json:"fps"`
	Fourcc    float64   `json:"fourcc"`
	Frames    int       `json:"frames"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewContent() *Content {
	c := Content{}
	return &c
}

func SelectContent(gId int) ([]*Content, error) {
	if gId == -1 {
		return Content{}.Order("id", "asc").All().Query()
	} else {
		return Content{}.Order("id", "asc").And("group_id", gId).Query()
	}
}

func SelectContentPager(gId int, page int) ([]*Content, error) {
	if gId == -1 {
		return Content{}.Order("id", "asc").Limit(ContentPageNum).Offset((page - 1) * ContentPageNum).Query()
	} else {
		return Content{}.Order("id", "asc").And("group_id", gId).Limit(ContentPageNum).Offset((page - 1) * ContentPageNum).Query()
	}
}
