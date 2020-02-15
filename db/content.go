package db

import (
	"time"
)

//+AR
type Content struct {
	ID        int       `json:"id" db:"pk"`
	GroupId   int       `json:"group_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	FPS       float64   `json:"fps"`
	Fourcc    float64   `json:"fourcc"`
	Frames    int       `json:"frames"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SelectContent(gId int) ([]*Content, error) {
	if gId == -1 {
		return Content{}.Order("id", "asc").All().Query()
	} else {
		return Content{}.Order("id", "asc").And("group_id", gId).Query()
	}
}
