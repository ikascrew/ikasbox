package ikasbox

import (
	"log"
	"os"

	mc "github.com/ikascrew/core/multicast"
	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"
	"github.com/ikascrew/ikasbox/handler"

	"golang.org/x/xerrors"
)

func Start(opts ...config.Option) error {

	err := config.Set(opts...)
	if err != nil {
		return xerrors.Errorf("config setting : %w", err)
	}

	conf := config.Get()

	switch conf.SubCommand {
	case "start":
		err = start()
	case "group":
		err = setGroup()
	case "project":
		err = setProject()
	case "init":
		err = create()
	}

	if err != nil {
		return xerrors.Errorf("subcommand error: %w", err)
	}

	return nil
}

func start() error {

	err := db.Open()
	if err != nil {
		return xerrors.Errorf("Database Open : %w", err)
	}

	go func() {
		c, err := mc.NewServer(
			mc.ServerName("ikasbox"),
			mc.Type(mc.TypeIkasbox),
		)
		if err != nil {
			log.Printf("multicast server error:%+v", err)
			return
		}
		err = c.Dial()
		if err != nil {
			log.Printf("multicast dial error:%+v", err)
			return
		}
	}()

	err = handler.Listen()
	if err != nil {
		return xerrors.Errorf("Handler Listen : %w", err)
	}

	return nil
}

func create() error {

	conf := config.Get()

	if _, err := os.Stat(conf.DatabasePath); err == nil {
		return xerrors.Errorf("the file already exists:%s", conf.DatabasePath)
	}

	err := db.Open()
	if err != nil {
		return xerrors.Errorf("db open error: %w", err)
	}

	err = db.CreateTables()
	if err != nil {
		return xerrors.Errorf("create tables error: %w", err)
	}

	return nil
}
