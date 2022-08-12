package db

import (
	"time"

	"golang.org/x/xerrors"
)

const CreateProjectGroupsSQL = `
CREATE TABLE [PROJECT_GROUPS] (
    [id] INTEGER PRIMARY KEY AUTOINCREMENT,
    [project_id] INTEGER,
    [group_id] INTEGER,
    [seq] INTEGER,
    [created_at] DATETIME,
    [updated_at] DATETIME
)
`

//+AR
type ProjectGroup struct {
	ID        int       `json:"id" db:"pk"`
	ProjectID int       `json:"project_id"`
	GroupID   int       `json:"group_id"`
	Seq       int       `json:"seq"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProjectGroup() *ProjectGroup {
	pg := ProjectGroup{}
	pg.ID = 0
	return &pg
}

func getProjectGroupList(id int) ([]*ProjectGroup, error) {
	pg := NewProjectGroup()
	pgs, err := pg.Where("project_id", id).Order("seq", "desc").Query()
	if err != nil {
		return nil, xerrors.Errorf("select project_groups error: %w", err)
	}
	return pgs, nil
}

func SelectProjectGroupList(id int) ([]*Group, error) {

	pgs, err := getProjectGroupList(id)
	if err != nil {
		return nil, xerrors.Errorf("getProjectGroupList() error: %w", err)
	}

	//TODO Join
	groups := make([]*Group, len(pgs))
	for idx, elm := range pgs {
		grp, err := FindGroup(elm.GroupID)
		if err != nil {
			return nil, xerrors.Errorf("select groups error: %w", err)
		}
		groups[idx] = grp
	}

	return groups, nil
}

func AddProjectGroup(pId, gId int) error {

	pgs, err := getProjectGroupList(pId)
	if err != nil {
		return xerrors.Errorf("getProjectGroupList() error: %w", err)
	}

	pg := NewProjectGroup()

	pg.ProjectID = pId
	pg.GroupID = gId
	pg.Seq = len(pgs) + 1
	pg.CreatedAt = time.Now()
	pg.UpdatedAt = time.Now()

	_, arerr := pg.Save(false)
	if arerr != nil {
		return xerrors.Errorf("content save: %w", arerr)
	}

	return nil
}

func SelectProjectContentList(id int) ([]*Content, error) {

	groups, err := SelectProjectGroupList(id)
	if err != nil {
		return nil, xerrors.Errorf("SelectProjectGroupList() error: %w", err)
	}

	contentList := make([]*Content, 0, 1000)

	for _, elm := range groups {
		list, err := SelectContent(elm.ID)
		if err != nil {
			return nil, xerrors.Errorf("select contents: %w", err)
		}
		contentList = append(contentList, list...)
	}

	return contentList, nil
}
