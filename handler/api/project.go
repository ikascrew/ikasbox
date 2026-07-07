package api

import (
	"github.com/ikascrew/ikasbox/db"
	"golang.org/x/xerrors"
)

type ProjectView struct {
	Paging db.Paging `json:"paging"`
}

type ProjectViewReturn struct {
	Projects []*db.Project `json:"projects"`
	Paging   db.Paging     `json:"paging"`
	Status
}

func newProjectView() Parameter {
	var pv ProjectView
	return &pv
}

func (pv *ProjectView) Processing() (Return, error) {

	var ret ProjectViewReturn

	projects, err := db.SelectPagingProject(&pv.Paging)
	if err != nil {
		return nil, xerrors.Errorf("select paging group: %w", err)
	}

	ret.Projects = projects
	ret.Paging = pv.Paging
	ret.success = true
	return &ret, nil
}

type ProjectGroupView struct {
	ProjectId int `json:"projectId"`
}

type ProjectGroupViewReturn struct {
	Project   *db.Project `json:"project"`
	Groups    []*db.Group `json:"groups"`
	AllGroups []*db.Group `json:"allGroups"`
	Status
}

func newProjectGroupView() Parameter {
	var pv ProjectGroupView
	return &pv
}

func (pv *ProjectGroupView) Processing() (Return, error) {

	var ret ProjectGroupViewReturn
	id := pv.ProjectId

	project, err := db.SelectProject(id)
	if err != nil {
		return nil, xerrors.Errorf("db.SelectProject() error: %w", err)
	}
	ret.Project = project

	groups, err := db.SelectProjectGroupList(id)
	if err != nil {
		return nil, xerrors.Errorf("select paging group: %w", err)
	}
	ret.Groups = groups

	all, err := db.SelectGroup()
	if err != nil {
		return nil, xerrors.Errorf("select group error: %w", err)
	}
	ret.AllGroups = all

	ret.success = true
	return &ret, nil
}

type ProjectGroupAdd struct {
	ProjectId int `json:"projectId"`
	GroupId   int `json:"groupId"`
}

type ProjectGroupAddReturn struct {
	Status
}

func newProjectGroupAdd() Parameter {
	var pa ProjectGroupAdd
	return &pa
}

func (pa *ProjectGroupAdd) Processing() (Return, error) {

	var ret ProjectGroupAddReturn

	err := db.AddProjectGroup(pa.ProjectId, pa.GroupId)
	if err != nil {
		return nil, xerrors.Errorf("AddProjectGroup() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}

type ProjectContentView struct {
	ProjectId int `json:"projectId"`
}

type ProjectContentViewReturn struct {
	Project  *db.Project   `json:"project"`
	Contents []*db.Content `json:"contents"`
	Status
}

func newProjectContentView() Parameter {
	var pcv ProjectContentView
	return &pcv
}

func (pcv *ProjectContentView) Processing() (Return, error) {

	var ret ProjectContentViewReturn

	project, err := db.SelectProject(pcv.ProjectId)
	if err != nil {
		return nil, xerrors.Errorf("db.SelectProject() error: %w", err)
	}
	ret.Project = project

	contents, err := db.SelectProjectContentList(pcv.ProjectId)
	if err != nil {
		return nil, xerrors.Errorf("db.SelectProjectContentList() error: %w", err)
	}
	ret.Contents = contents

	ret.success = true
	return &ret, nil
}

type ProjectRegister struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ProjectRegisterReturn struct {
	Status
}

func newProjectRegister() Parameter {
	var gr ProjectRegister
	return &gr
}

func (pr *ProjectRegister) Processing() (Return, error) {

	var ret ProjectRegisterReturn

	err := db.RegisterProject(pr.Name, pr.Width, pr.Height)
	if err != nil {
		return nil, xerrors.Errorf("RegisterProject() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}

type ProjectRename struct {
	ProjectId int    `json:"projectId"`
	Name      string `json:"name"`
}

type ProjectRenameReturn struct {
	Status
}

func newProjectRename() Parameter {
	var pr ProjectRename
	return &pr
}

func (pr *ProjectRename) Processing() (Return, error) {

	var ret ProjectRenameReturn

	err := db.RenameProject(pr.ProjectId, pr.Name)
	if err != nil {
		return nil, xerrors.Errorf("db.RenameProject() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}

type ProjectDelete struct {
	ProjectId int `json:"projectId"`
}

type ProjectDeleteReturn struct {
	Status
}

func newProjectDelete() Parameter {
	var pd ProjectDelete
	return &pd
}

func (pd *ProjectDelete) Processing() (Return, error) {

	var ret ProjectDeleteReturn

	err := db.DeleteProject(pd.ProjectId)
	if err != nil {
		return nil, xerrors.Errorf("db.DeleteProject() error: %w", err)
	}

	ret.success = true
	return &ret, nil
}
