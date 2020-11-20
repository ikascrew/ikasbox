package ikasbox

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

func setProject() error {

	var err error
	conf := config.Get()
	args := conf.Arguments

	switch conf.Function {
	case "list", "add":
		var projects []*db.Project
		projects, err = viewProjects()

		if conf.Function == "list" {
			break
		}

		pId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("project id error: %w", err)
		}

		for _, elm := range projects {
			if elm.ID == pId {
				err = addContent(pId, args[1:]...)
				break
			}
		}

	case "register":
		err = registerProject(args[0])
	default:
		err = fmt.Errorf("sub command error[%s]", conf.Function)
	}

	if err != nil {
		return xerrors.Errorf("sub command(%s) error: %w", conf.Function, err)
	}

	return nil
}

func viewProjects() ([]*db.Project, error) {

	projects, err := db.SelectProject()
	if err != nil {
		return nil, xerrors.Errorf("select project: %w", err)
	}

	for _, elm := range projects {
		fmt.Printf("[%d]:%s\n", elm.ID, elm.Name)
	}

	return projects, nil

}

func registerProject(name string) error {

	project := db.Project{
		Name:      name,
		Width:     1280,
		Height:    720,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := project.Save(true)
	if err != nil {
		return xerrors.Errorf("register project: %w", err)
	}

	fmt.Printf("New Project:%s[%d]\n", name, project.ID)

	return nil
}

func addContent(pId int, groups ...string) error {

	for _, buf := range groups {
		gId, err := strconv.Atoi(buf)
		if err != nil {
			return fmt.Errorf("project id error(%s): %w", buf, err)
		}

		err = addProjectContent(gId, pId)
		if err != nil {
			return xerrors.Errorf("register project error(%d): %w", gId, err)
		}
	}

	return nil
}

func addProjectContent(gId, pId int) error {

	var err error

	contentList, err := db.SelectContent(gId)
	if err != nil {
		return xerrors.Errorf("select content: %w", err)
	}

	for _, elm := range contentList {
		content := db.ProjectContent{
			ProjectID: pId,
			ContentID: elm.ID,
			Type:      "file",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, arerr := content.Save(false)
		if arerr != nil {
			return xerrors.Errorf("content save: %w", err)
		}
	}

	return nil
}
