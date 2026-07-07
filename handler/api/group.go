package api

import (
	"log"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/contentimport"
	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

type GroupView struct {
	Paging db.Paging `json:"paging"`
}

type GroupViewReturn struct {
	Groups []*db.Group `json:"groups"`
	Paging db.Paging   `json:"paging"`
	Status
}

func newGroupView() Parameter {
	var gv GroupView
	return &gv
}

func (gv *GroupView) Processing() (Return, error) {

	var ret GroupViewReturn

	groups, err := db.SelectPagingGroup(&gv.Paging)
	if err != nil {
		return nil, xerrors.Errorf("select paging group: %w", err)
	}

	ret.Groups = groups
	ret.Paging = gv.Paging
	ret.success = true
	return &ret, nil
}

type GroupRegister struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type GroupRegisterReturn struct {
	GroupId int `json:"groupId"`
	Status
}

func newGroupRegister() Parameter {
	var gr GroupRegister
	return &gr
}

func (gr *GroupRegister) Processing() (Return, error) {

	var ret GroupRegisterReturn

	groupId, err := db.RegisterGroup(gr.Name, gr.Path)
	if err != nil {
		return nil, xerrors.Errorf("RegisterGroup() error: %w", err)
	}
	ret.GroupId = groupId

	if gr.Path != "" {
		path := gr.Path
		extensions := config.Get().Extensions

		// Import runs in the background; the request returns immediately
		// without waiting for the (potentially slow) thumbnail generation.
		go func() {
			count, err := contentimport.ImportDirectory(groupId, path, extensions)
			if err != nil {
				log.Printf("group import error: %+v", err)
				return
			}
			log.Printf("group import completed: group[%d] path[%s] imported[%d]", groupId, path, count)
		}()
	}

	ret.success = true
	return &ret, nil
}

type GroupCheck struct {
	GroupId int `json:"groupId"`
}

type GroupCheckReturn struct {
	Missing []*db.Content `json:"missing"`
	Total   int           `json:"total"`
	Status
}

func newGroupCheck() Parameter {
	var gc GroupCheck
	return &gc
}

func (gc *GroupCheck) Processing() (Return, error) {

	var ret GroupCheckReturn

	contents, err := db.SelectContent(gc.GroupId)
	if err != nil {
		return nil, xerrors.Errorf("db.SelectContent() error: %w", err)
	}
	ret.Total = len(contents)

	missing, err := contentimport.CheckMissing(gc.GroupId)
	if err != nil {
		return nil, xerrors.Errorf("contentimport.CheckMissing() error: %w", err)
	}
	ret.Missing = missing

	ret.success = true
	return &ret, nil
}

type GroupRename struct {
	GroupId int    `json:"groupId"`
	Name    string `json:"name"`
}

type GroupRenameReturn struct {
	Status
}

func newGroupRename() Parameter {
	var gr GroupRename
	return &gr
}

func (gr *GroupRename) Processing() (Return, error) {

	var ret GroupRenameReturn

	err := db.RenameGroup(gr.GroupId, gr.Name)
	if err != nil {
		return nil, xerrors.Errorf("db.RenameGroup() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}

type GroupDelete struct {
	GroupId int `json:"groupId"`
}

type GroupDeleteReturn struct {
	Status
}

func newGroupDelete() Parameter {
	var gd GroupDelete
	return &gd
}

func (gd *GroupDelete) Processing() (Return, error) {

	var ret GroupDeleteReturn

	err := db.DeleteGroup(gd.GroupId)
	if err != nil {
		return nil, xerrors.Errorf("db.DeleteGroup() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}
