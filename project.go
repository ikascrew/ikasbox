package ikasbox

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ikascrew/core/util"
	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

func setProject() error {

	err := db.Open()
	if err != nil {
		return xerrors.Errorf("Database Open : %w", err)
	}

	conf := config.Get()
	args := conf.Arguments

	switch conf.Function {
	case "list", "add":

		var projects []*db.Project
		projects, err = viewProjects()
		if conf.Function == "list" {
			break
		}

		err = inputProject(projects)

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

	projects, err := db.SelectProjectList()
	if err != nil {
		return nil, xerrors.Errorf("select project: %w", err)
	}

	for _, elm := range projects {
		fmt.Printf("[%d]:%s\n", elm.ID, elm.Name)
	}

	return projects, nil

}

func registerProject(name string) error {

	var err error
	width := 1280
	height := 720

	fmt.Print("Width[1280]:")
	in := util.Input()

	if in != "" {
		width, err = strconv.Atoi(in)
		if err != nil {
			return xerrors.Errorf("input width: %w", err)
		}
	}

	fmt.Print("Height[720]:")
	in = util.Input()

	if in != "" {
		height, err = strconv.Atoi(in)
		if err != nil {
			return xerrors.Errorf("input height: %w", err)
		}
	}

	project := db.Project{
		Name:      name,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, arerr := project.Save(true)
	if arerr != nil {
		return xerrors.Errorf("register project: %w", arerr)
	}

	fmt.Printf("New Project:%s[%d]\n", name, project.ID)

	return nil
}

func inputProject(projects []*db.Project) error {

	fmt.Print("Select Project ID:")
	buf := util.Input()

	pId, err := strconv.Atoi(buf)
	if err != nil {
		return fmt.Errorf("project id error: %w", err)
	}

	current := -1
	for _, p := range projects {
		if p.ID == pId {
			current = pId
			break
		}
	}

	if current == -1 {
		return fmt.Errorf("not found project id[%d]", pId)
	}

	groups, err := viewGroups()
	if err != nil {
		return xerrors.Errorf("view groups: %w", err)
	}

	//TODO :現在のGroupを表示

	fmt.Print("Select Group:")
	buf = util.Input()

	gId, err := strconv.Atoi(buf)
	if err != nil {
		return fmt.Errorf("project id error: %w", err)
	}

	current = -1
	for _, g := range groups {
		if g.ID == gId {
			current = gId
			break
		}
	}

	if current == -1 {
		return fmt.Errorf("not found group id[%d]", gId)
	}

	err = db.AddProjectGroup(pId, gId)
	if err != nil {
		return xerrors.Errorf("add group: %w", err)
	}

	return nil
}
