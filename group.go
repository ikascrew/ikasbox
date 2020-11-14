package ikasbox

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
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

	conf := config.Get()
	var err error
	switch conf.Function {
	case "register":
		err = fmt.Errorf("not implemented.")
	case "import":
		err = importContent(conf.Arguments[0])
	case "check":
		err = check()
	case "list":
		err = fmt.Errorf("not implemented.")
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

func importContent(p string) error {

	conf := config.Get()

	//グループの一覧を表示
	group, err := ChooseGroup()
	if err != nil {
		return err
	}

	//ファイルの検索
	files, err := own.SearchDirectory(p, conf.Extensions)
	if err != nil {
		return err
	}

	//ファイルのソート
	own.SortFiles(files)
	if len(files) <= 0 {
		return fmt.Errorf("file not found[%s]", p)
	}
	fmt.Printf("[%s] target files[%d]. Register?[Y/n]:", p, len(files))

	in := input()
	if in != "Y" {
		return nil
	}

	g, err := db.FindGroup(group)
	if err != nil {
		return err
	}

	now := time.Now()
	_, arErr := g.Update(db.GroupParams{Path: p, UpdatedAt: now})
	if arErr != nil {
		return err
	}

	bar := pb.StartNew(len(files)).Prefix("Register Content")

	//work := r.Path + string(os.PathSeparator) + ".ikabox"

	thumb := path.Join("public", "images", "thumb")
	os.MkdirAll(thumb, 0777)

	for _, elm := range files {
		err := registerFile(thumb, group, elm)
		if err != nil {
			fmt.Printf("Error Register[%s][%s]\n", elm, err)
		}
		bar.Increment()
	}

	bar.FinishPrint("Register Content Completion")
	return nil
}

func registerFile(dir string, id int, f string) error {

	v, err := util.NewVideo(f)

	if err != nil {
		return err
	}

	frames := float64(v.Frames)
	images := make([]*gocv.Mat, 17)
	//半分の位置を取得
	m, err := v.GetImage(frames / 2.0)
	if err != nil {
		return err
	}

	images[0] = m

	for idx := 1; idx <= 16; idx++ {
		i, err := v.GetImage(frames / 16.0 * float64(idx))
		if err != nil {
			return err
		}
		images[idx] = i
	}

	db.Transaction(func(tx *sql.Tx) error {

		now := time.Now()
		//コンテンツを登録
		c := db.Content{
			GroupId:   id,
			Name:      filepath.Base(f),
			Path:      f,
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
			return err
		}

		//画像をリサイズ
		//r, err := util.ResizeImage(*m, 256, 144)
		//if err != nil {
		//return err
		//}
		//defer r.Close()

		bufId := strconv.FormatInt(int64(c.ID), 10)

		for idx, img := range images {

			ext := fmt.Sprintf("_%d.jpg", idx)
			if idx == 0 {
				ext = ".jpg"
			}

			//グループのパスに設定
			thumb := dir + string(os.PathSeparator) + bufId + ext

			if !img.Empty() {
				//サムネイルを作成
				err = util.WriteImage(thumb, *img)
				if err != nil {
					return err
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

	return err
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

	groups, err := db.SelectGroup()
	if err != nil {
		return -1, err
	}

	if len(groups) == 0 {
		return -1, fmt.Errorf("Group not exists")
	}

	for idx, elm := range groups {
		fmt.Printf("[%d] %s\n", idx+1, elm.Name)
	}
	fmt.Printf("Select[ %d - %d ] :", 1, len(groups))
	in := input()

	idx, err := strconv.Atoi(in)
	if err != nil {
		return -1, err
	}

	if idx <= 0 || idx > len(groups) {
		return -1, fmt.Errorf("Error Index value.")
	}

	fmt.Printf("Register group[%s]\n", groups[idx-1].Name)
	return groups[idx-1].ID, nil
}

func input() string {
	std := bufio.NewScanner(os.Stdin)
	std.Scan()
	text := std.Text()
	return text
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
		if ans := input(); ans == "Y" {
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
