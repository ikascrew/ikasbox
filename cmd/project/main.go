package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

const registerProjectNumber = -1

func main() {

	//go run main.go list
	//go run main.go register "name" {groups ...}
	//go run main.go add (project id) {groups ...}

	flag.Parse()
	args := flag.Args()

	err := run(args)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Println("Success")
}

func run(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("project command arguments requied[list,register,add]")
	}

	sub := args[0]
	var err error

	switch sub {
	case "list":
		//project list
	case "register":
		if len(args) < 3 {
			return fmt.Errorf("register command arguments projectname , groups ...")
		}
		err = addContent(registerProjectNumber, args[1:]...)
	case "add":
		if len(args) < 3 {
			return fmt.Errorf("add command arguments projectid , groups ...")
		}
		pId, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("project id error: %w", err)
		}
		err = addContent(pId, args[2:]...)
	default:
		err = fmt.Errorf("sub command error[%s]", sub)
	}

	if err != nil {
		return xerrors.Errorf("sub command(%s) error: %w", sub, err)
	}

	return nil
}

func addContent(pId int, groups ...string) error {

	if pId == registerProjectNumber {
		//register project
		name := groups[0]
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

		groups = groups[1:]
		pId = project.ID

		fmt.Printf("New Project:%s[%d]\n", name, pId)
	}

	for _, buf := range groups {
		gId, err := strconv.Atoi(buf)
		if err != nil {
			return fmt.Errorf("project id error(%s): %w", buf, err)
		}

		err = registerProject(gId, pId)
		if err != nil {
			return xerrors.Errorf("register project error(%d): %w", gId, err)
		}
	}

	return nil
}

func registerProject(gId, pId int) error {

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
