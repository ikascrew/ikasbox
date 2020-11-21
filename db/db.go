package db

import (
	"database/sql"
	"log"

	"github.com/ikascrew/ikasbox/config"
	"golang.org/x/xerrors"

	_ "github.com/mattn/go-sqlite3"
)

func Open() error {

	c := config.Get()
	dbfile := c.DatabasePath

	log.Println("open database:" + dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return xerrors.Errorf("open database error: %w", err)
	}
	Use(db)

	LogMode(c.Debug)
	return nil
}

func Transaction(fn func(tx *sql.Tx) error) (err error) {

	var tx *sql.Tx

	tx, err = db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		rec := recover()
		if rec != nil {
			err = xerrors.Errorf("tx recover error: %w", rec)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = fn(tx)
	return
}

func CreateTables() error {

	_, err := db.Exec(CreateGroupsSQL)
	if err != nil {
		return xerrors.Errorf("create groups: %w", err)
	}
	_, err = db.Exec(CreateContentsSQL)
	if err != nil {
		return xerrors.Errorf("create contents: %w", err)
	}
	_, err = db.Exec(CreateContentThumbnailsSQL)
	if err != nil {
		return xerrors.Errorf("create content_thumbnails: %w", err)
	}
	_, err = db.Exec(CreateProjectsSQL)
	if err != nil {
		return xerrors.Errorf("create projects: %w", err)
	}

	_, err = db.Exec(CreateProjectContentsSQL)
	if err != nil {
		return xerrors.Errorf("create project_contents: %w", err)
	}
	return nil
}
