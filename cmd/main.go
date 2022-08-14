package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/ikasbox"
	"github.com/ikascrew/ikasbox/config"
	"golang.org/x/xerrors"
)

//ビルド時に指定するDB名
// -ldflags "-X main.embedDB=xxxx.db"
var embedDB string

var dbfile *string
var exts *string

func init() {
	dbfile = flag.String("db", "ikasbox.db", "sqlite database file")
	exts = flag.String("ext", "*.mp4,*.mpeg,*.png,*.jpg,*.jpeg", "import extention")
}

func main() {

	flag.Parse()
	args := flag.Args()

	err := run(args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Println("success")
}

func run(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("ikasbox subcommand arguments required.(start,group,project,init)")
	}

	opts := createOptions(args)
	err := ikasbox.Start(opts...)

	if err != nil {
		return xerrors.Errorf("ikasbox command error: %w", err)
	}

	return nil
}

func createOptions(args []string) []config.Option {

	opts := make([]config.Option, 0)

	//ビルド時に埋め込みがある場合にそちらを優先
	db := *dbfile
	if embedDB != "" {
		db = embedDB
	}

	//DBファイル設定
	opts = append(opts, config.Path(db))
	opts = append(opts, config.Extension(*exts))
	//引数設定
	opts = append(opts, config.Argument(args))

	return opts
}
