package main

import (
	"bufio"
	"database/sql"
	"flag"
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
	"gopkg.in/cheggaaa/pb.v1"
)

var dbfile *string

func init() {
	dbfile = flag.String("db", "ikasbox.db", "sqlite database file")
}

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		log.Println("argument error")
		os.Exit(1)
	}

	err := config.Set(config.Path(*dbfile))
	if err != nil {
		panic(err)
	}

	//引数でインポートする
	reg := Register{
		Path: args[0],
		Ext:  []string{"*.mp4", "*.mpeg", "*.png", "*.jpg", "*.jpeg"},
		Blob: false,
	}

	err = RegisterContent(reg)
	if err != nil {
		panic(err)
	}
}

type Register struct {
	Path string
	Ext  []string
	Blob bool
}

func RegisterContent(r Register) error {

	//グループの一覧を表示
	group, err := ChooseGroup()
	if err != nil {
		return err
	}

	//ファイルの検索
	files, err := own.SearchDirectory(r.Path, r.Ext)
	if err != nil {
		return err
	}

	//ファイルのソート
	own.SortFiles(files)
	if len(files) <= 0 {
		return nil
	}
	fmt.Printf("[%s] target files[%d]. Register?[Y/n]:", r.Path, len(files))

	in := Input()
	if in != "Y" {
		return nil
	}

	g, err := db.FindGroup(group)
	if err != nil {
		return err
	}

	now := time.Now()
	_, arErr := g.Update(db.GroupParams{Path: r.Path, UpdatedAt: now})
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
	in := Input()

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

func Input() string {
	std := bufio.NewScanner(os.Stdin)
	std.Scan()
	text := std.Text()
	return text
}
