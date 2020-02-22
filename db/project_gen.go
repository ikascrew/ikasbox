// generated by argen; DO NOT EDIT
package db

import (
	"fmt"

	"github.com/monochromegane/argen"
)

type ProjectRelation struct {
	src *Project
	*ar.Relation
}

func (m *Project) newRelation() *ProjectRelation {
	r := &ProjectRelation{
		m,
		ar.NewRelation(db, logger).Table("projects"),
	}
	r.Select(
		"id",
		"name",
		"width",
		"height",
		"default_content",
		"created_at",
		"updated_at",
	)

	return r
}

func (m Project) Select(columns ...string) *ProjectRelation {
	return m.newRelation().Select(columns...)
}

func (r *ProjectRelation) Select(columns ...string) *ProjectRelation {
	cs := []string{}
	for _, c := range columns {
		if r.src.isColumnName(c) {
			cs = append(cs, fmt.Sprintf("projects.%s", c))
		} else {
			cs = append(cs, c)
		}
	}
	r.Relation.Columns(cs...)
	return r
}

func (m Project) Find(id int) (*Project, error) {
	return m.newRelation().Find(id)
}

func (r *ProjectRelation) Find(id int) (*Project, error) {
	return r.FindBy("id", id)
}

func (m Project) FindBy(cond string, args ...interface{}) (*Project, error) {
	return m.newRelation().FindBy(cond, args...)
}

func (r *ProjectRelation) FindBy(cond string, args ...interface{}) (*Project, error) {
	return r.Where(cond, args...).Limit(1).QueryRow()
}

func (m Project) First() (*Project, error) {
	return m.newRelation().First()
}

func (r *ProjectRelation) First() (*Project, error) {
	return r.Order("id", "ASC").Limit(1).QueryRow()
}

func (m Project) Last() (*Project, error) {
	return m.newRelation().Last()
}

func (r *ProjectRelation) Last() (*Project, error) {
	return r.Order("id", "DESC").Limit(1).QueryRow()
}

func (m Project) Where(cond string, args ...interface{}) *ProjectRelation {
	return m.newRelation().Where(cond, args...)
}

func (r *ProjectRelation) Where(cond string, args ...interface{}) *ProjectRelation {
	r.Relation.Where(cond, args...)
	return r
}

func (r *ProjectRelation) And(cond string, args ...interface{}) *ProjectRelation {
	r.Relation.And(cond, args...)
	return r
}

func (m Project) Order(column, order string) *ProjectRelation {
	return m.newRelation().Order(column, order)
}

func (r *ProjectRelation) Order(column, order string) *ProjectRelation {
	r.Relation.OrderBy(column, order)
	return r
}

func (m Project) Limit(limit int) *ProjectRelation {
	return m.newRelation().Limit(limit)
}

func (r *ProjectRelation) Limit(limit int) *ProjectRelation {
	r.Relation.Limit(limit)
	return r
}

func (m Project) Offset(offset int) *ProjectRelation {
	return m.newRelation().Offset(offset)
}

func (r *ProjectRelation) Offset(offset int) *ProjectRelation {
	r.Relation.Offset(offset)
	return r
}

func (m Project) Group(group string, groups ...string) *ProjectRelation {
	return m.newRelation().Group(group, groups...)
}

func (r *ProjectRelation) Group(group string, groups ...string) *ProjectRelation {
	r.Relation.GroupBy(group, groups...)
	return r
}

func (r *ProjectRelation) Having(cond string, args ...interface{}) *ProjectRelation {
	r.Relation.Having(cond, args...)
	return r
}

func (m Project) IsValid() (bool, *ar.Errors) {
	result := true
	errors := &ar.Errors{}
	var on ar.On
	if m.IsNewRecord() {
		on = ar.OnCreate()
	} else {
		on = ar.OnUpdate()
	}
	rules := map[string]*ar.Validation{}
	for name, rule := range rules {
		if ok, errs := ar.NewValidator(rule).On(on).IsValid(m.fieldValueByName(name)); !ok {
			result = false
			errors.SetErrors(name, errs)
		}
	}
	customs := []*ar.Validation{}
	for _, rule := range customs {
		custom := ar.NewValidator(rule).On(on).Custom()
		custom(errors)
	}
	if len(errors.Messages) > 0 {
		result = false
	}
	return result, errors
}

type ProjectParams Project

func (m Project) Build(p ProjectParams) *Project {
	return &Project{
		ID:             p.ID,
		Name:           p.Name,
		Width:          p.Width,
		Height:         p.Height,
		DefaultContent: p.DefaultContent,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func (m Project) Create(p ProjectParams) (*Project, *ar.Errors) {
	n := m.Build(p)
	_, errs := n.Save()
	return n, errs
}

func (m *Project) IsNewRecord() bool {
	return ar.IsZero(m.ID)
}

func (m *Project) IsPersistent() bool {
	return !m.IsNewRecord()
}

func (m *Project) Save(validate ...bool) (bool, *ar.Errors) {
	if len(validate) == 0 || len(validate) > 0 && validate[0] {
		if ok, errs := m.IsValid(); !ok {
			return false, errs
		}
	}
	errs := &ar.Errors{}
	if m.IsNewRecord() {
		ins := ar.NewInsert(db, logger).Table("projects").Params(map[string]interface{}{
			"name":            m.Name,
			"width":           m.Width,
			"height":          m.Height,
			"default_content": m.DefaultContent,
			"created_at":      m.CreatedAt,
			"updated_at":      m.UpdatedAt,
		})

		if result, err := ins.Exec(); err != nil {
			errs.AddError("base", err)
			return false, errs
		} else {
			if lastId, err := result.LastInsertId(); err == nil {
				m.ID = int(lastId)
			}
		}
		return true, nil
	} else {
		upd := ar.NewUpdate(db, logger).Table("projects").Params(map[string]interface{}{
			"id":              m.ID,
			"name":            m.Name,
			"width":           m.Width,
			"height":          m.Height,
			"default_content": m.DefaultContent,
			"created_at":      m.CreatedAt,
			"updated_at":      m.UpdatedAt,
		}).Where("id", m.ID)

		if _, err := upd.Exec(); err != nil {
			errs.AddError("base", err)
			return false, errs
		}
		return true, nil
	}
}

func (m *Project) Update(p ProjectParams) (bool, *ar.Errors) {

	if !ar.IsZero(p.ID) {
		m.ID = p.ID
	}
	if !ar.IsZero(p.Name) {
		m.Name = p.Name
	}
	if !ar.IsZero(p.Width) {
		m.Width = p.Width
	}
	if !ar.IsZero(p.Height) {
		m.Height = p.Height
	}
	if !ar.IsZero(p.DefaultContent) {
		m.DefaultContent = p.DefaultContent
	}
	if !ar.IsZero(p.CreatedAt) {
		m.CreatedAt = p.CreatedAt
	}
	if !ar.IsZero(p.UpdatedAt) {
		m.UpdatedAt = p.UpdatedAt
	}
	return m.Save()
}

func (m *Project) UpdateColumns(p ProjectParams) (bool, *ar.Errors) {

	if !ar.IsZero(p.ID) {
		m.ID = p.ID
	}
	if !ar.IsZero(p.Name) {
		m.Name = p.Name
	}
	if !ar.IsZero(p.Width) {
		m.Width = p.Width
	}
	if !ar.IsZero(p.Height) {
		m.Height = p.Height
	}
	if !ar.IsZero(p.DefaultContent) {
		m.DefaultContent = p.DefaultContent
	}
	if !ar.IsZero(p.CreatedAt) {
		m.CreatedAt = p.CreatedAt
	}
	if !ar.IsZero(p.UpdatedAt) {
		m.UpdatedAt = p.UpdatedAt
	}
	return m.Save(false)
}

func (m *Project) Destroy() (bool, *ar.Errors) {
	return m.Delete()
}

func (m *Project) Delete() (bool, *ar.Errors) {
	errs := &ar.Errors{}
	if _, err := ar.NewDelete(db, logger).Table("projects").Where("id", m.ID).Exec(); err != nil {
		errs.AddError("base", err)
		return false, errs
	}
	return true, nil
}

func (m Project) DeleteAll() (bool, *ar.Errors) {
	errs := &ar.Errors{}
	if _, err := ar.NewDelete(db, logger).Table("projects").Exec(); err != nil {
		errs.AddError("base", err)
		return false, errs
	}
	return true, nil
}

func (r *ProjectRelation) Query() ([]*Project, error) {
	rows, err := r.Relation.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*Project{}
	for rows.Next() {
		row := &Project{}
		err := rows.Scan(row.fieldPtrsByName(r.Relation.GetColumns())...)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (r *ProjectRelation) QueryRow() (*Project, error) {
	row := &Project{}
	err := r.Relation.QueryRow(row.fieldPtrsByName(r.Relation.GetColumns())...)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (m Project) Exists() bool {
	return m.newRelation().Exists()
}

func (m Project) Count(column ...string) int {
	return m.newRelation().Count(column...)
}

func (m Project) All() *ProjectRelation {
	return m.newRelation().All()
}

func (r *ProjectRelation) All() *ProjectRelation {
	return r
}

func (m *Project) fieldValueByName(name string) interface{} {
	switch name {
	case "id", "projects.id":
		return m.ID
	case "name", "projects.name":
		return m.Name
	case "width", "projects.width":
		return m.Width
	case "height", "projects.height":
		return m.Height
	case "default_content", "projects.default_content":
		return m.DefaultContent
	case "created_at", "projects.created_at":
		return m.CreatedAt
	case "updated_at", "projects.updated_at":
		return m.UpdatedAt
	default:
		return ""
	}
}

func (m *Project) fieldPtrByName(name string) interface{} {
	switch name {
	case "id", "projects.id":
		return &m.ID
	case "name", "projects.name":
		return &m.Name
	case "width", "projects.width":
		return &m.Width
	case "height", "projects.height":
		return &m.Height
	case "default_content", "projects.default_content":
		return &m.DefaultContent
	case "created_at", "projects.created_at":
		return &m.CreatedAt
	case "updated_at", "projects.updated_at":
		return &m.UpdatedAt
	default:
		return nil
	}
}

func (m *Project) fieldPtrsByName(names []string) []interface{} {
	fields := []interface{}{}
	for _, n := range names {
		f := m.fieldPtrByName(n)
		fields = append(fields, f)
	}
	return fields
}

func (m *Project) isColumnName(name string) bool {
	for _, c := range m.columnNames() {
		if c == name {
			return true
		}
	}
	return false
}

func (m *Project) columnNames() []string {
	return []string{
		"id",
		"name",
		"width",
		"height",
		"default_content",
		"created_at",
		"updated_at",
	}
}