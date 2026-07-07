package db

import (
	"database/sql"
	"fmt"

	"github.com/ikascrew/ikasbox/config"
	"golang.org/x/xerrors"

	_ "github.com/mattn/go-sqlite3"
)

func Open() error {

	c := config.Get()
	dbfile := c.DatabasePath

	fmt.Println("open database:" + dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return xerrors.Errorf("open database error: %w", err)
	}
	Use(db)

	LogMode(c.Debug)

	if err := migrate(); err != nil {
		return xerrors.Errorf("migrate error: %w", err)
	}
	return nil
}

// migrate は既存データベースへのスキーマ変更を適用する。
// 冪等になるよう、列の存在を確認してから ALTER TABLE する
func migrate() error {

	// init 前(contents 自体が無い)は何もしない
	var name string
	err := db.QueryRow(
		`SELECT name FROM sqlite_master WHERE type='table' AND name='contents'`).Scan(&name)
	if err != nil {
		return nil
	}

	rows, err := db.Query(`PRAGMA table_info(contents)`)
	if err != nil {
		return xerrors.Errorf("table_info(contents): %w", err)
	}
	defer rows.Close()

	has := false
	for rows.Next() {
		var cid int
		var colName, typ string
		var notnull int
		var dflt interface{}
		var pk int
		if err := rows.Scan(&cid, &colName, &typ, &notnull, &dflt, &pk); err != nil {
			return xerrors.Errorf("table_info scan: %w", err)
		}
		if colName == "params" {
			has = true
		}
	}
	if err := rows.Err(); err != nil {
		return xerrors.Errorf("table_info rows: %w", err)
	}

	if !has {
		if _, err := db.Exec(`ALTER TABLE contents ADD COLUMN params TEXT`); err != nil {
			return xerrors.Errorf("add column params: %w", err)
		}
		fmt.Println("migrate: add column contents.params")
	}
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

	fmt.Println("Create Groups")
	_, err := db.Exec(CreateGroupsSQL)
	if err != nil {
		return xerrors.Errorf("create groups: %w", err)
	}

	fmt.Println("Create Contents")
	_, err = db.Exec(CreateContentsSQL)
	if err != nil {
		return xerrors.Errorf("create contents: %w", err)
	}

	fmt.Println("Create ContentThumbnails")
	_, err = db.Exec(CreateContentThumbnailsSQL)
	if err != nil {
		return xerrors.Errorf("create content_thumbnails: %w", err)
	}

	fmt.Println("Create Projects")
	_, err = db.Exec(CreateProjectsSQL)
	if err != nil {
		return xerrors.Errorf("create projects: %w", err)
	}

	fmt.Println("Create ProjectGroups")
	_, err = db.Exec(CreateProjectGroupsSQL)
	if err != nil {
		return xerrors.Errorf("create project_groups: %w", err)
	}
	return nil
}
