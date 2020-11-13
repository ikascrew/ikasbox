package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

func main() {

	err := run()
	if err != nil {
		log.Println("%+v", err)
		os.Exit(1)
	}

	log.Println("Success")
}

func run() error {
	//コンテンツの全件取得
	//パスにコンテンツがあるか？
	contents, err := getContents()
	if err != nil {
		return xerrors.Errorf("content all: %w", err)
	}

	//プログレスバーを表示
	nothings, err := checkContent(contents)
	if err != nil {
		return xerrors.Errorf("check: %w", err)
	}

	if len(nothings) > 0 {
		for _, con := range nothings {
			fmt.Printf("%d:%s(%s)\n", con.ID, con.Name, con.Path)
		}

		log.Printf("nothing %d/all %d Delete?[Y/n]:", len(nothings), len(contents))

		//delete?
		if ans := input(); ans == "Y" {
		}

	} else {
		log.Println("exists all")
	}

	return nil
}

func getContents() ([]*db.Content, error) {
	contents, err := db.SelectContent(-1)
	if err != nil {
		return nil, xerrors.Errorf("select content error: %w", err)
	}

	return contents, nil
}

func checkContent(all []*db.Content) ([]*db.Content, error) {

	nothings := make([]*db.Content, 0, len(all))
	for _, content := range all {
		if _, err := os.Stat(content.Path); err != nil {
			nothings = append(nothings, content)
		}
	}

	return nothings, nil
}

func input() string {
	std := bufio.NewScanner(os.Stdin)
	std.Scan()
	text := std.Text()
	return text
}
