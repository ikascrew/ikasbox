package db

import (
	"strconv"
	"time"

	ar "github.com/monochromegane/argen"
)

const CreateProjectContentsSQL = `
CREATE TABLE [PROJECT_CONTENTS] (
    [id] INTEGER PRIMARY KEY AUTOINCREMENT,
    [project_id] INTEGER NOT NULL,
    [content_id] INTEGER NOT NULL,
    [type] VARCHAR(256) NOT NULL,
    [created_at] DATETIME,
    [updated_at] DATETIME
)
`

//+AR
type ProjectContent struct {
	ID        int       `json:"id" db:"pk"`
	ProjectID int       `json:"project_id" db:"fk"`
	ContentID int       `json:"content_id" db:"fk"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p ProjectContent) belongsToContent() *ar.Association {
	return nil
}

func SelectProjectContentList(project string) (*Project, []*ProjectContent, error) {

	id, err := strconv.Atoi(project)
	if err != nil {
		return nil, nil, err
	}

	p, arerr := Project{}.Find(id)
	if arerr != nil {
		return nil, nil, arerr
	}

	contentList, arerr := ProjectContent{}.JoinsContent().And("project_id", id).Query()
	if arerr != nil {
		return nil, nil, arerr
	}

	return p, contentList, nil
}
