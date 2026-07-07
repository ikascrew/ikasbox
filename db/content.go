package db

import (
	"time"

	"golang.org/x/xerrors"
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
    [params] TEXT,
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
	Params    string    `json:"params"`
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

func SelectPagingContent(gId int, pg *Paging) ([]*Content, error) {

	var dao Content

	if gId == -1 {
		cnt := dao.All().Count()
		pg.SetCount(cnt)
		return dao.Order("id", "asc").Limit(pg.Limit).Offset((pg.Current - 1) * pg.Limit).Query()
	}

	cnt := dao.Where("group_id", gId).Count()
	pg.SetCount(cnt)

	return dao.Where("group_id", gId).Order("id", "asc").Limit(pg.Limit).Offset((pg.Current - 1) * pg.Limit).Query()
}

func DeleteContent(id int) error {

	c, err := Content{}.Find(id)
	if err != nil {
		return xerrors.Errorf("content find: %w", err)
	}

	if err := DeleteContentThumbnails(c.ID); err != nil {
		return xerrors.Errorf("content thumbnail delete: %w", err)
	}

	if _, arErr := c.Delete(); arErr != nil {
		return xerrors.Errorf("content delete: %w", arErr)
	}

	return nil
}
