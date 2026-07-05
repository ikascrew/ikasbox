// Package contentimport scans a directory for media files and registers
// them as contents (with generated thumbnails) for a group. It is shared
// by the CLI ("ikasbox group import") and the web API so both drive the
// same import/thumbnail-generation logic.
package contentimport

import (
	"bytes"
	"database/sql"
	"image/jpeg"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikascrew/core/util"
	"github.com/ikascrew/ikasbox/db"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

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

	//TODO FileかImageかを拡張子でやっていると思うのでだめ
	v, err := util.NewVideo(f)
	if err != nil {
		return xerrors.Errorf("load video[%s]: %w", f, err)
	}
	defer v.Close()

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
			return xerrors.Errorf("content save: %w", arErr)
		}

		for idx, img := range images {

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
