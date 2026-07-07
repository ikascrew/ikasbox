package db

import (
	"time"

	"golang.org/x/xerrors"
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

func NewGroup() *Group {
	g := Group{}
	return &g
}

func SelectGroup() ([]*Group, error) {

	groups, err := Group{}.Order("id", "asc").All().Query()
	if err != nil {
		return nil, xerrors.Errorf("select group: %w", err)
	}
	return groups, nil
}

func SelectPagingGroup(pg *Paging) ([]*Group, error) {

	var grp Group
	cnt := grp.Count()
	pg.SetCount(cnt)

	offset := (pg.Current - 1) * pg.Limit

	groups, err := Group{}.Order("updated_at", "desc").Limit(pg.Limit).Offset(offset).Query()
	if err != nil {
		return nil, xerrors.Errorf("select group: %w", err)
	}
	return groups, nil
}

func FindGroup(id int) (*Group, error) {
	g, err := Group{}.Find(id)
	if err != nil {
		return nil, xerrors.Errorf("find group: %w", err)
	}
	return g, nil
}

func RegisterGroup(name string, path string) (int, error) {

	now := time.Now()
	g := Group{
		Name:      name,
		Path:      path,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := g.Save(false)
	if err != nil {
		return 0, xerrors.Errorf("Group Save() error: %w", err)
	}

	return g.ID, nil
}

// DeleteGroup removes a group together with all of its contents and
// their thumbnails.
func DeleteGroup(id int) error {

	g, err := Group{}.Find(id)
	if err != nil {
		return xerrors.Errorf("group find: %w", err)
	}

	contents, err := SelectContent(id)
	if err != nil {
		return xerrors.Errorf("select content: %w", err)
	}

	for _, c := range contents {
		if err := DeleteContent(c.ID); err != nil {
			return xerrors.Errorf("delete content[%d]: %w", c.ID, err)
		}
	}

	if _, arErr := g.Delete(); arErr != nil {
		return xerrors.Errorf("group delete: %w", arErr)
	}

	return nil
}
