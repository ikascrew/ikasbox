package db

import (
	"time"

	"golang.org/x/xerrors"
)

const CreateProjectGroupsSQL = `
CREATE TABLE [PROJECT_GROUPS] (
    [ID] INTEGER PRIMARY KEY AUTOINCREMENT,
    [PROJECT_ID] INTEGER,
    [GROUP_ID] INTEGER,
    [CREATED_AT] DATETIME,
    [UPDATED_AT] DATETIME
)
`

//+AR
type ProjectGroup struct {
	ID        int       `json:"id" db:"pk"`
	ProjectID int       `json:"project_id"`
	GroupID   int       `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProjectGroup() *ProjectGroup {
	pg := ProjectGroup{}
	return &pg
}

func SelectProjectContentList(id int) ([]*Content, error) {

	pg := NewProjectGroup()

	groups, err := pg.Where("project_id", id).Order("created_at", "asc").Query()
	if err != nil {
		return nil, xerrors.Errorf("select project_groups: %w", err)
	}

	contentList := make([]*Content, 0, 1000)

	for _, elm := range groups {
		list, err := SelectContent(elm.GroupID)
		if err != nil {
			return nil, xerrors.Errorf("select contents: %w", err)
		}
		contentList = append(contentList, list...)
	}

	return contentList, nil
}
