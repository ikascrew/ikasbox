package api

import (
	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

type ContentView struct {
	GroupId int       `json:"groupId"`
	Paging  db.Paging `json:"paging"`
}

type ContentViewReturn struct {
	Group    *db.Group     `json:"group"`
	Contents []*db.Content `json:"contents"`
	Paging   db.Paging     `json:"paging"`
	Status
}

func newContentView() Parameter {
	var cv ContentView
	return &cv
}

func (cv *ContentView) Processing() (Return, error) {

	var ret ContentViewReturn

	var g *db.Group
	if cv.GroupId == -1 {
		g = &db.Group{ID: -1, Name: "すべてのコンテンツ"}
	} else {
		found, err := db.FindGroup(cv.GroupId)
		if err != nil {
			return nil, xerrors.Errorf("db.FindGroup() error: %w", err)
		}
		g = found
	}

	contents, err := db.SelectPagingContent(cv.GroupId, &cv.Paging)
	if err != nil {
		return nil, xerrors.Errorf("select paging group: %w", err)
	}

	ret.Group = g
	ret.Contents = contents
	ret.Paging = cv.Paging
	ret.success = true

	return &ret, nil
}

type ContentFind struct {
	Id int `json:"id"`
}

type ContentFindReturn struct {
	Content *db.Content `json:"content"`
	Status
}

func newContentFind() Parameter {
	var cf ContentFind
	return &cf
}

func (cf *ContentFind) Processing() (Return, error) {

	var ret ContentFindReturn

	content, err := db.Content{}.Find(cf.Id)
	if err != nil {
		return nil, xerrors.Errorf("db.Content{}.Find() error: %w", err)
	}

	ret.Content = content
	ret.success = true

	return &ret, nil
}
