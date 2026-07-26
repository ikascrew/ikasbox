// Package contentimport scans a directory for media files and registers
// them as contents (with generated thumbnails) for a group. It is shared
// by the CLI ("ikasbox group import") and the web API so both drive the
// same import/thumbnail-generation logic.
package contentimport

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikascrew/core/util"
	"github.com/ikascrew/ikasbox/db"
	"github.com/ikascrew/plugin/video"

	"golang.org/x/xerrors"
)

// サムネイルの枚数(seq 0 は中間フレーム、1..16 がタイムライン)
const thumbnailNum = 17

// SearchFiles finds and sorts the files under path matching extensions.
func SearchFiles(path string, extensions []string) ([]string, error) {

	files, err := util.SearchDirectory(path, extensions)
	if err != nil {
		return nil, xerrors.Errorf("search directory: %w", err)
	}

	util.SortFiles(files)
	return files, nil
}

// ImportDirectory searches path for media files, updates the group's
// registered path, and registers each file as a content. It returns the
// number of files successfully registered; files that fail to register are
// logged and skipped.
func ImportDirectory(groupId int, path string, extensions []string) (int, error) {

	files, err := SearchFiles(path, extensions)
	if err != nil {
		return 0, xerrors.Errorf("search files: %w", err)
	}

	if len(files) <= 0 {
		return 0, xerrors.Errorf("file not found[%s]", path)
	}

	g, err := db.FindGroup(groupId)
	if err != nil {
		return 0, xerrors.Errorf("find group: %w", err)
	}

	if _, arErr := g.Update(db.GroupParams{Path: path, UpdatedAt: time.Now()}); arErr != nil {
		return 0, xerrors.Errorf("group update: %w", arErr)
	}

	count := 0
	for _, f := range files {
		if err := RegisterFile(groupId, f); err != nil {
			log.Println(err)
			continue
		}
		count++
	}

	return count, nil
}

// RegisterFile loads a single media file, generates thumbnails for it, and
// saves it as a content of the given group.
func RegisterFile(id int, f string) error {

	v, err := util.NewVideo(f)
	if err != nil {
		return xerrors.Errorf("load video[%s]: %w", f, err)
	}
	defer v.Close()

	// 型は plugin の正語彙(file/img)で保存する
	typ := "file"
	if isImage(f) {
		typ = "img"
	}

	frames := float64(v.Frames)
	images := make([]image.Image, thumbnailNum)

	if isImage(f) {
		// 静止画は位置によらず同一フレームなので1回だけ変換して共有する
		// (image.Image は不変なので Mat と違い共有しても安全)
		m, err := v.GetGoImage(0)
		if err != nil {
			return xerrors.Errorf("get image: %w", err)
		}
		for idx := range images {
			images[idx] = m
		}
	} else {
		//半分の位置を取得
		m, err := v.GetGoImage(frames / 2.0)
		if err != nil {
			return xerrors.Errorf("get image(root): %w", err)
		}

		images[0] = m
		for idx := 0; idx <= 15; idx++ {
			i, err := v.GetGoImage(frames/16.0*float64(idx) + 1)
			if err != nil {
				return xerrors.Errorf("get image(%d): %w", idx, err)
			}
			images[idx+1] = i
		}
	}

	now := time.Now()
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

	if err := saveContentWithThumbnails(&c, images); err != nil {
		return xerrors.Errorf("register content: %w", err)
	}

	return nil
}

// RegisterGenerated は実体ファイルを持たない生成型コンテンツ
// (cd, terminal など)を登録する。params は JSON 文字列で、
// 妥当性はプラグインを実際に生成して検証し、サムネイルも
// プラグインにフレームを描画させて作る
func RegisterGenerated(groupId int, name, typ, params string) (*db.Content, error) {

	if _, err := db.FindGroup(groupId); err != nil {
		return nil, xerrors.Errorf("find group[%d]: %w", groupId, err)
	}

	t := video.Normalize(typ)

	// params の検証(プラグインの実生成)とフレーム描画は plugin 側に
	// 閉じている(壊れた JSON を本番まで持ち込まない)
	images, err := video.RenderThumbnails(t, params, thumbnailNum)
	if err != nil {
		return nil, xerrors.Errorf("render thumbnails[%s]: %w", t, err)
	}

	bounds := images[0].Bounds()

	now := time.Now()
	c := db.Content{
		GroupId:   groupId,
		Name:      name,
		Path:      "",
		Type:      t,
		Params:    params,
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		FPS:       30,
		Frames:    0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := saveContentWithThumbnails(&c, images); err != nil {
		return nil, xerrors.Errorf("register content: %w", err)
	}

	return &c, nil
}

// saveContentWithThumbnails はコンテンツとサムネイル一式を
// 1トランザクションで保存する。
// argen 生成の Save()/Insert() はパッケージレベルの *sql.DB 直結で
// tx を通せないため、原子性が必要なこの保存だけは tx へ直接 INSERT する
func saveContentWithThumbnails(c *db.Content, images []image.Image) error {

	err := db.Transaction(func(tx *sql.Tx) error {

		res, err := tx.Exec(
			`insert into contents(group_id,name,type,path,params,width,height,fps,fourcc,frames,created_at,updated_at)
			 values (?,?,?,?,?,?,?,?,?,?,?,?)`,
			c.GroupId, c.Name, c.Type, c.Path, c.Params,
			c.Width, c.Height, c.FPS, c.Fourcc, c.Frames,
			c.CreatedAt, c.UpdatedAt)
		if err != nil {
			return xerrors.Errorf("content insert: %w", err)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return xerrors.Errorf("content last insert id: %w", err)
		}
		c.ID = int(lastId)

		for idx, img := range images {

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, img, nil)
			if err != nil {
				return xerrors.Errorf("convert image: %w", err)
			}

			_, err = tx.Exec(
				"insert into content_thumbnails(id,seq,data) values (?,?,?)",
				c.ID, idx, buf.Bytes())
			if err != nil {
				return xerrors.Errorf("thumbnail insert: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		// ロールバック時に採番済み ID が残らないようにする
		c.ID = 0
	}
	return err
}

// CheckMissing returns the contents whose file no longer exists on disk.
// Pass groupId of -1 to check every content across all groups.
func CheckMissing(groupId int) ([]*db.Content, error) {

	contents, err := db.SelectContent(groupId)
	if err != nil {
		return nil, xerrors.Errorf("select content: %w", err)
	}

	missing := make([]*db.Content, 0, len(contents))
	for _, c := range contents {
		if _, err := os.Stat(c.Path); err != nil {
			missing = append(missing, c)
		}
	}

	return missing, nil
}

// CheckDuplicates returns groups of contents that share the same file path,
// i.e. the same file registered more than once. Each returned slice has 2 or
// more entries. Pass groupId of -1 to check every content across all groups.
func CheckDuplicates(groupId int) ([][]*db.Content, error) {

	contents, err := db.SelectContent(groupId)
	if err != nil {
		return nil, xerrors.Errorf("select content: %w", err)
	}

	byPath := make(map[string][]*db.Content)
	for _, c := range contents {
		byPath[c.Path] = append(byPath[c.Path], c)
	}

	duplicates := make([][]*db.Content, 0)
	for _, group := range byPath {
		if len(group) > 1 {
			duplicates = append(duplicates, group)
		}
	}

	return duplicates, nil
}

func isImage(f string) bool {
	if strings.Index(f, ".jpg") != -1 ||
		strings.Index(f, ".jpeg") != -1 ||
		strings.Index(f, ".png") != -1 {
		return true
	}
	return false
}
