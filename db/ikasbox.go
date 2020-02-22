package db

import (
	"database/sql"
	"fmt"

	"github.com/ikascrew/ikasbox/config"

	_ "github.com/mattn/go-sqlite3"
)

func init() {

	c := config.Get()
	dbfile := c.DatabasePath

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	Use(db)

	LogMode(c.Debug)
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
			err = fmt.Errorf("Panic=", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = fn(tx)
	return
}
