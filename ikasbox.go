package ikasbox

import (
	"log"

	mc "github.com/ikascrew/core/multicast"
	"github.com/ikascrew/ikasbox/config"
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
	}

	if err != nil {
		return xerrors.Errorf("subcommand error: %w", err)
	}

	return nil
}

func start() error {

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

	err := handler.Listen()
	if err != nil {
		return xerrors.Errorf("Handler Listen : %w", err)
	}

	return nil
}
