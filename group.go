package ikasbox

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ikascrew/core/util"
	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/contentimport"
	"github.com/ikascrew/ikasbox/db"

	"golang.org/x/xerrors"
	"gopkg.in/cheggaaa/pb.v1"
)

func setGroup() error {

	err := db.Open()
	if err != nil {
		return xerrors.Errorf("Database Open : %w", err)
	}

	conf := config.Get()

	switch conf.Function {
	case "register":
		err = registerGroup(conf.Arguments[0])
	case "import":
		err = importContent(conf.Arguments[0])
	case "check":
		err = check()
	case "list":
		_, err = viewGroups()
	case "remove":
		err = removeGroup(conf.Arguments[0])
	default:
		err = fmt.Errorf("not found function")
	}

	if err != nil {
		return xerrors.Errorf("function error[%s]: %w", conf.Function, err)
	}

	return nil
}

func viewGroups() ([]*db.Group, error) {
	groups, err := db.SelectGroup()
	if err != nil {
		return nil, xerrors.Errorf("select group: %w", err)
	}
	for _, elm := range groups {
		fmt.Printf("[%d] %s\n", elm.ID, elm.Name)
	}
	return groups, nil
}

func registerGroup(name string) error {

	group := db.Group{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := group.Save(true)
	if err != nil {
		return xerrors.Errorf("register group: %w", err)
	}

	fmt.Printf("New Group:%s[%d]\n", name, group.ID)

	return nil
}

func removeGroup(arg string) error {

	id, err := strconv.Atoi(arg)
	if err != nil {
		return xerrors.Errorf("input id: %w", err)
	}

	fmt.Printf("Delete group[%d] and all its contents?[Y/n]:", id)
	if ans := util.Input(); ans != "Y" {
		return nil
	}

	if err := db.DeleteGroup(id); err != nil {
		return xerrors.Errorf("delete group: %w", err)
	}

	fmt.Printf("Deleted Group[%d]\n", id)
	return nil
}

func importContent(p string) error {

	conf := config.Get()

	//グループの一覧を表示
	id, err := ChooseGroup()
	if err != nil {
		return xerrors.Errorf("choose group: %w", err)
	}

	//ファイルの検索
	files, err := contentimport.SearchFiles(p, conf.Extensions)
	if err != nil {
		return xerrors.Errorf("search files: %w", err)
	}

	if len(files) <= 0 {
		return fmt.Errorf("file not found[%s]", p)
	}
	fmt.Printf("[%s] target files[%d]. Register?[Y/n]:", p, len(files))

	in := util.Input()
	if in != "Y" {
		return nil
	}

	g, err := db.FindGroup(id)
	if err != nil {
		return xerrors.Errorf("find group: %w", err)
	}

	now := time.Now()
	_, arErr := g.Update(db.GroupParams{Path: p, UpdatedAt: now})
	if arErr != nil {
		return xerrors.Errorf("group update: %w", err)
	}

	bar := pb.StartNew(len(files)).Prefix("Register Content")

	//wg := &sync.WaitGroup{}
	//sem := make(chan struct{}, 10)
	//defer close(sem)

	for _, elm := range files {
		//sem <- struct{}{}
		//wg.Add(1)
		//go func(name string) {
		//defer func() {
		//wg.Done()
		//<-sem
		//}()

		err := contentimport.RegisterFile(id, elm)
		if err != nil {
			log.Println(err)
		}
		bar.Increment()
		//}(elm)
	}

	//wg.Wait()

	bar.FinishPrint("Register Content Completion")
	return nil
}

func ChooseGroup() (int, error) {

	groups, err := viewGroups()
	if err != nil {
		return -1, xerrors.Errorf("view Groups: %w", err)
	}

	groupMap := make(map[int]*db.Group)
	for _, elm := range groups {
		groupMap[elm.ID] = elm
	}

	fmt.Printf("Select GroupID :")
	in := util.Input()

	id, err := strconv.Atoi(in)
	if err != nil {
		return -1, xerrors.Errorf("input id error: %w", err)
	}

	if g, ok := groupMap[id]; ok {
		fmt.Printf("Register group[%s]\n", g.Name)
		return g.ID, nil
	}
	return -1, fmt.Errorf("Error ID[%d]", id)

}

func check() error {
	//コンテンツの全件取得、パスにコンテンツがあるか？
	contents, err := db.SelectContent(-1)
	if err != nil {
		return xerrors.Errorf("select content error: %w", err)
	}

	nothings, err := contentimport.CheckMissing(-1)
	if err != nil {
		return xerrors.Errorf("check: %w", err)
	}

	if len(nothings) > 0 {
		for _, con := range nothings {
			fmt.Printf("%d:%s(%s)\n", con.ID, con.Name, con.Path)
		}

		log.Printf("nothing %d/all %d Delete?[Y/n]:", len(nothings), len(contents))

		//delete?
		if ans := util.Input(); ans == "Y" {
			return fmt.Errorf("not implemented")
		}

	} else {
		log.Println("exists all")
	}

	return nil
}
