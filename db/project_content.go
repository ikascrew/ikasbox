package db

import (
	"strconv"
	"time"

	"github.com/monochromegane/argen"
)

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
