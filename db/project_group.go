package db

import (
	"time"

	"golang.org/x/xerrors"
)

const CreateProjectGroupsSQL = `
CREATE TABLE [PROJECT_GROUPS] (
    [PROJECT_ID] INTEGER,
    [GROUP_ID] INTEGER,
    [CREATED_AT] DATETIME,
    [UPDATED_AT] DATETIME,
  PRIMARY KEY([PROJECT_ID],[GROUP_ID])
)
`

//+AR
type ProjectGroup struct {
	ID        int       `json:"project_id"`
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

	groups, err := pg.Where("id", id).Order("create_at", "asc").Query()
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
