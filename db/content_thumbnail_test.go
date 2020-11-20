package db_test

import (
	"database/sql"
	"io/ioutil"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func TestThumbnail(t *testing.T) {

	f, err := sql.Open("sqlite3", "thumbnail.db")
	if err != nil {
		t.Errorf("initialize db error: %+v", err)
	}
	db.Use(f)

	th := db.ContentThumbnail{}
	th.ID = 1
	th.Seq = 1

	bytes, err := ioutil.ReadFile("sample.jpg")
	if err != nil {
		t.Errorf("read file error: %+v", err)
	}

	th.Data = bytes

	err = th.Insert()
	if err != nil {
		t.Errorf("insert error: %+v", err)
	}
}

func TestSelectContentThumbnails(t *testing.T) {

	f, err := sql.Open("sqlite3", "thumbnail.db")
	if err != nil {
		t.Errorf("initialize db error: %+v", err)
	}
	db.Use(f)

	ths, err := db.SelectContentThumbnails(1)
	if err != nil {
		t.Errorf("select content thumbnails error: %+v", err)
	}

	if len(ths) != 1 {
		t.Errorf("select content thumbnail count error")
	}

	err = ioutil.WriteFile("thumbnail.jpg", ths[0].Data, 0744)
	if err != nil {
		t.Errorf("create file count error: %+v", err)
	}

}
