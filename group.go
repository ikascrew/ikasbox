package ikasbox

import (
	"bytes"
	"database/sql"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ikascrew/core/util"
	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"
	own "github.com/ikascrew/ikasbox/util"
	"gocv.io/x/gocv"
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
		err = fmt.Errorf("not implemented.")
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

func importContent(p string) error {

	conf := config.Get()

	//グループの一覧を表示
	id, err := ChooseGroup()
	if err != nil {
		return xerrors.Errorf("choose group: %w", err)
	}

	//ファイルの検索
	files, err := own.SearchDirectory(p, conf.Extensions)
	if err != nil {
		return xerrors.Errorf("search directory: %w", err)
	}

	//ファイルのソート
	own.SortFiles(files)
	if len(files) <= 0 {
		return fmt.Errorf("file not found[%s]", p)
	}
	fmt.Printf("[%s] target files[%d]. Register?[Y/n]:", p, len(files))

	in := own.Input()
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

	for _, elm := range files {
		err := registerFile(id, elm)
		if err != nil {
			return xerrors.Errorf("Error Register[%s]: %w", elm, err)
		}
		bar.Increment()
	}

	bar.FinishPrint("Register Content Completion")
	return nil
}

func registerFile(id int, f string) error {

	//TODO FileかImageかを拡張子でやっていると思うのでだめ
	v, err := util.NewVideo(f)
	if err != nil {
		return xerrors.Errorf("load video[%s]: %w", f, err)
	}

	typ := "file"
	if isImage(f) {
		typ = "image"
	}

	frames := float64(v.Frames)
	images := make([]*gocv.Mat, 17)
	//半分の位置を取得
	m, err := v.GetImage(frames / 2.0)
	if err != nil {
		return xerrors.Errorf("get image(root): %w", err)
	}

	images[0] = m
	for idx := 0; idx <= 15; idx++ {
		i, err := v.GetImage(frames/16.0*float64(idx) + 1)
		if err != nil {
			return xerrors.Errorf("get image(%d): %w", idx, err)
		}
		images[idx+1] = i
	}

	err = db.Transaction(func(tx *sql.Tx) error {

		now := time.Now()
		//コンテンツを登録
		c := db.Content{
			GroupId:   id,
			Name:      filepath.Base(f),
			Path:      f,
			Type:      typ,
			Width:     v.Width,
			Height:    v.Height,
			FPS:       v.FPS,
			Fourcc:    v.FOURCC,
			Frames:    v.Frames,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, arErr := c.Save()
		if arErr != nil {
			return xerrors.Errorf("content save: %w", err)
		}

		for idx, img := range images {
			if !img.Empty() {

				thumb := db.ContentThumbnail{}

				thumb.ID = c.ID
				thumb.Seq = idx
				goimg, err := img.ToImage()
				if err != nil {
					return xerrors.Errorf("mat to image: %w", err)
				}

				buf := new(bytes.Buffer)
				err = jpeg.Encode(buf, goimg, nil)
				if err != nil {
					return xerrors.Errorf("convert image: %w", err)
				}

				thumb.Data = buf.Bytes()

				err = thumb.Insert()
				if err != nil {
					return xerrors.Errorf("thumbnail insert: %w", err)
				}

				if !isImage(f) {
					img.Close()
				}
			}
		}

		if isImage(f) {
			m.Close()
		}
		return nil
	})

	if err != nil {
		return xerrors.Errorf("register content: %w", err)
	}

	return nil
}

func isImage(f string) bool {
	if strings.Index(f, ".jpg") != -1 ||
		strings.Index(f, ".jpeg") != -1 ||
		strings.Index(f, ".png") != -1 {
		return true
	}
	return false
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
	in := own.Input()

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
	//コンテンツの全件取得
	//パスにコンテンツがあるか？
	contents, err := getContents()
	if err != nil {
		return xerrors.Errorf("content all: %w", err)
	}

	//プログレスバーを表示
	nothings, err := checkContent(contents)
	if err != nil {
		return xerrors.Errorf("check: %w", err)
	}

	if len(nothings) > 0 {
		for _, con := range nothings {
			fmt.Printf("%d:%s(%s)\n", con.ID, con.Name, con.Path)
		}

		log.Printf("nothing %d/all %d Delete?[Y/n]:", len(nothings), len(contents))

		//delete?
		if ans := own.Input(); ans == "Y" {
			return fmt.Errorf("not implemented")
		}

	} else {
		log.Println("exists all")
	}

	return nil
}

func getContents() ([]*db.Content, error) {
	contents, err := db.SelectContent(-1)
	if err != nil {
		return nil, xerrors.Errorf("select content error: %w", err)
	}

	return contents, nil
}

func checkContent(all []*db.Content) ([]*db.Content, error) {

	nothings := make([]*db.Content, 0, len(all))
	for _, content := range all {
		if _, err := os.Stat(content.Path); err != nil {
			nothings = append(nothings, content)
		}
	}

	return nothings, nil
}
