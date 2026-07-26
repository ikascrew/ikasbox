package api

import (
	"strings"

	"github.com/ikascrew/ikasbox/contentimport"
	"github.com/ikascrew/ikasbox/db"

	"github.com/ikascrew/plugin/video"
	"github.com/ikascrew/plugin/video/param"

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

// ContentSpec は登録フォームを組み立てるための型一覧を返す。
// フィールド定義はプラグインの自己申告(video.Spec)をそのまま流すだけで、
// ikasbox は中身を解釈しない
type ContentSpec struct{}

// ContentSpecEntry は1つの型の登録フォーム定義。
// Fields が空(未知の型)の場合、UI は生 JSON 入力へフォールバックする
type ContentSpecEntry struct {
	Type string `json:"type"`
	// Generative は実体ファイルを持たない生成型か。UI が
	// 「ディレクトリ import」と「params 登録」を出し分けるための情報
	Generative bool          `json:"generative"`
	Fields     []param.Field `json:"fields"`
}

type ContentSpecReturn struct {
	Types []ContentSpecEntry `json:"types"`
	Status
}

func newContentSpec() Parameter {
	var cs ContentSpec
	return &cs
}

func (cs *ContentSpec) Processing() (Return, error) {

	var ret ContentSpecReturn

	for _, t := range video.Types() {
		ret.Types = append(ret.Types, ContentSpecEntry{
			Type:       t,
			Generative: video.IsGenerative(t),
			Fields:     video.Spec(t),
		})
	}

	ret.success = true

	return &ret, nil
}

// ContentRegister は params を指定してコンテンツを1件登録する。
// params の妥当性検証とサムネイル生成は contentimport 側
// (実際にプラグインを生成する)に任せる
type ContentRegister struct {
	GroupId int    `json:"groupId"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Params  string `json:"params"`
}

type ContentRegisterReturn struct {
	Content *db.Content `json:"content"`
	Status
}

func newContentRegister() Parameter {
	var cr ContentRegister
	return &cr
}

func (cr *ContentRegister) Processing() (Return, error) {

	var ret ContentRegisterReturn

	if cr.GroupId <= 0 {
		return nil, xerrors.Errorf("groupId is required")
	}
	if strings.TrimSpace(cr.Name) == "" {
		return nil, xerrors.Errorf("name is required")
	}

	content, err := contentimport.RegisterGenerated(cr.GroupId, cr.Name, cr.Type, cr.Params)
	if err != nil {
		return nil, xerrors.Errorf("RegisterGenerated() error: %w", err)
	}

	ret.Content = content
	ret.success = true

	return &ret, nil
}
