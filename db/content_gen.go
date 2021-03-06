// generated by argen; DO NOT EDIT
package db

import (
	"fmt"

	"github.com/monochromegane/argen"
)

type ContentRelation struct {
	src *Content
	*ar.Relation
}

func (m *Content) newRelation() *ContentRelation {
	r := &ContentRelation{
		m,
		ar.NewRelation(db, logger).Table("contents"),
	}
	r.Select(
		"id",
		"group_id",
		"name",
		"type",
		"path",
		"width",
		"height",
		"fps",
		"fourcc",
		"frames",
		"created_at",
		"updated_at",
	)

	return r
}

func (m Content) Select(columns ...string) *ContentRelation {
	return m.newRelation().Select(columns...)
}

func (r *ContentRelation) Select(columns ...string) *ContentRelation {
	cs := []string{}
	for _, c := range columns {
		if r.src.isColumnName(c) {
			cs = append(cs, fmt.Sprintf("contents.%s", c))
		} else {
			cs = append(cs, c)
		}
	}
	r.Relation.Columns(cs...)
	return r
}

func (m Content) Find(id int) (*Content, error) {
	return m.newRelation().Find(id)
}

func (r *ContentRelation) Find(id int) (*Content, error) {
	return r.FindBy("id", id)
}

func (m Content) FindBy(cond string, args ...interface{}) (*Content, error) {
	return m.newRelation().FindBy(cond, args...)
}

func (r *ContentRelation) FindBy(cond string, args ...interface{}) (*Content, error) {
	return r.Where(cond, args...).Limit(1).QueryRow()
}

func (m Content) First() (*Content, error) {
	return m.newRelation().First()
}

func (r *ContentRelation) First() (*Content, error) {
	return r.Order("id", "ASC").Limit(1).QueryRow()
}

func (m Content) Last() (*Content, error) {
	return m.newRelation().Last()
}

func (r *ContentRelation) Last() (*Content, error) {
	return r.Order("id", "DESC").Limit(1).QueryRow()
}

func (m Content) Where(cond string, args ...interface{}) *ContentRelation {
	return m.newRelation().Where(cond, args...)
}

func (r *ContentRelation) Where(cond string, args ...interface{}) *ContentRelation {
	r.Relation.Where(cond, args...)
	return r
}

func (r *ContentRelation) And(cond string, args ...interface{}) *ContentRelation {
	r.Relation.And(cond, args...)
	return r
}

func (m Content) Order(column, order string) *ContentRelation {
	return m.newRelation().Order(column, order)
}

func (r *ContentRelation) Order(column, order string) *ContentRelation {
	r.Relation.OrderBy(column, order)
	return r
}

func (m Content) Limit(limit int) *ContentRelation {
	return m.newRelation().Limit(limit)
}

func (r *ContentRelation) Limit(limit int) *ContentRelation {
	r.Relation.Limit(limit)
	return r
}

func (m Content) Offset(offset int) *ContentRelation {
	return m.newRelation().Offset(offset)
}

func (r *ContentRelation) Offset(offset int) *ContentRelation {
	r.Relation.Offset(offset)
	return r
}

func (m Content) Group(group string, groups ...string) *ContentRelation {
	return m.newRelation().Group(group, groups...)
}

func (r *ContentRelation) Group(group string, groups ...string) *ContentRelation {
	r.Relation.GroupBy(group, groups...)
	return r
}

func (r *ContentRelation) Having(cond string, args ...interface{}) *ContentRelation {
	r.Relation.Having(cond, args...)
	return r
}

func (m Content) IsValid() (bool, *ar.Errors) {
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

type ContentParams Content

func (m Content) Build(p ContentParams) *Content {
	return &Content{
		ID:        p.ID,
		GroupId:   p.GroupId,
		Name:      p.Name,
		Type:      p.Type,
		Path:      p.Path,
		Width:     p.Width,
		Height:    p.Height,
		FPS:       p.FPS,
		Fourcc:    p.Fourcc,
		Frames:    p.Frames,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (m Content) Create(p ContentParams) (*Content, *ar.Errors) {
	n := m.Build(p)
	_, errs := n.Save()
	return n, errs
}

func (m *Content) IsNewRecord() bool {
	return ar.IsZero(m.ID)
}

func (m *Content) IsPersistent() bool {
	return !m.IsNewRecord()
}

func (m *Content) Save(validate ...bool) (bool, *ar.Errors) {
	if len(validate) == 0 || len(validate) > 0 && validate[0] {
		if ok, errs := m.IsValid(); !ok {
			return false, errs
		}
	}
	errs := &ar.Errors{}
	if m.IsNewRecord() {
		ins := ar.NewInsert(db, logger).Table("contents").Params(map[string]interface{}{
			"group_id":   m.GroupId,
			"name":       m.Name,
			"type":       m.Type,
			"path":       m.Path,
			"width":      m.Width,
			"height":     m.Height,
			"fps":        m.FPS,
			"fourcc":     m.Fourcc,
			"frames":     m.Frames,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
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
		upd := ar.NewUpdate(db, logger).Table("contents").Params(map[string]interface{}{
			"id":         m.ID,
			"group_id":   m.GroupId,
			"name":       m.Name,
			"type":       m.Type,
			"path":       m.Path,
			"width":      m.Width,
			"height":     m.Height,
			"fps":        m.FPS,
			"fourcc":     m.Fourcc,
			"frames":     m.Frames,
			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
		}).Where("id", m.ID)

		if _, err := upd.Exec(); err != nil {
			errs.AddError("base", err)
			return false, errs
		}
		return true, nil
	}
}

func (m *Content) Update(p ContentParams) (bool, *ar.Errors) {

	if !ar.IsZero(p.ID) {
		m.ID = p.ID
	}
	if !ar.IsZero(p.GroupId) {
		m.GroupId = p.GroupId
	}
	if !ar.IsZero(p.Name) {
		m.Name = p.Name
	}
	if !ar.IsZero(p.Type) {
		m.Type = p.Type
	}
	if !ar.IsZero(p.Path) {
		m.Path = p.Path
	}
	if !ar.IsZero(p.Width) {
		m.Width = p.Width
	}
	if !ar.IsZero(p.Height) {
		m.Height = p.Height
	}
	if !ar.IsZero(p.FPS) {
		m.FPS = p.FPS
	}
	if !ar.IsZero(p.Fourcc) {
		m.Fourcc = p.Fourcc
	}
	if !ar.IsZero(p.Frames) {
		m.Frames = p.Frames
	}
	if !ar.IsZero(p.CreatedAt) {
		m.CreatedAt = p.CreatedAt
	}
	if !ar.IsZero(p.UpdatedAt) {
		m.UpdatedAt = p.UpdatedAt
	}
	return m.Save()
}

func (m *Content) UpdateColumns(p ContentParams) (bool, *ar.Errors) {

	if !ar.IsZero(p.ID) {
		m.ID = p.ID
	}
	if !ar.IsZero(p.GroupId) {
		m.GroupId = p.GroupId
	}
	if !ar.IsZero(p.Name) {
		m.Name = p.Name
	}
	if !ar.IsZero(p.Type) {
		m.Type = p.Type
	}
	if !ar.IsZero(p.Path) {
		m.Path = p.Path
	}
	if !ar.IsZero(p.Width) {
		m.Width = p.Width
	}
	if !ar.IsZero(p.Height) {
		m.Height = p.Height
	}
	if !ar.IsZero(p.FPS) {
		m.FPS = p.FPS
	}
	if !ar.IsZero(p.Fourcc) {
		m.Fourcc = p.Fourcc
	}
	if !ar.IsZero(p.Frames) {
		m.Frames = p.Frames
	}
	if !ar.IsZero(p.CreatedAt) {
		m.CreatedAt = p.CreatedAt
	}
	if !ar.IsZero(p.UpdatedAt) {
		m.UpdatedAt = p.UpdatedAt
	}
	return m.Save(false)
}

func (m *Content) Destroy() (bool, *ar.Errors) {
	return m.Delete()
}

func (m *Content) Delete() (bool, *ar.Errors) {
	errs := &ar.Errors{}
	if _, err := ar.NewDelete(db, logger).Table("contents").Where("id", m.ID).Exec(); err != nil {
		errs.AddError("base", err)
		return false, errs
	}
	return true, nil
}

func (m Content) DeleteAll() (bool, *ar.Errors) {
	errs := &ar.Errors{}
	if _, err := ar.NewDelete(db, logger).Table("contents").Exec(); err != nil {
		errs.AddError("base", err)
		return false, errs
	}
	return true, nil
}

func (r *ContentRelation) Query() ([]*Content, error) {
	rows, err := r.Relation.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*Content{}
	for rows.Next() {
		row := &Content{}
		err := rows.Scan(row.fieldPtrsByName(r.Relation.GetColumns())...)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (r *ContentRelation) QueryRow() (*Content, error) {
	row := &Content{}
	err := r.Relation.QueryRow(row.fieldPtrsByName(r.Relation.GetColumns())...)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (m Content) Exists() bool {
	return m.newRelation().Exists()
}

func (m Content) Count(column ...string) int {
	return m.newRelation().Count(column...)
}

func (m Content) All() *ContentRelation {
	return m.newRelation().All()
}

func (r *ContentRelation) All() *ContentRelation {
	return r
}

func (m *Content) fieldValueByName(name string) interface{} {
	switch name {
	case "id", "contents.id":
		return m.ID
	case "group_id", "contents.group_id":
		return m.GroupId
	case "name", "contents.name":
		return m.Name
	case "type", "contents.type":
		return m.Type
	case "path", "contents.path":
		return m.Path
	case "width", "contents.width":
		return m.Width
	case "height", "contents.height":
		return m.Height
	case "fps", "contents.fps":
		return m.FPS
	case "fourcc", "contents.fourcc":
		return m.Fourcc
	case "frames", "contents.frames":
		return m.Frames
	case "created_at", "contents.created_at":
		return m.CreatedAt
	case "updated_at", "contents.updated_at":
		return m.UpdatedAt
	default:
		return ""
	}
}

func (m *Content) fieldPtrByName(name string) interface{} {
	switch name {
	case "id", "contents.id":
		return &m.ID
	case "group_id", "contents.group_id":
		return &m.GroupId
	case "name", "contents.name":
		return &m.Name
	case "type", "contents.type":
		return &m.Type
	case "path", "contents.path":
		return &m.Path
	case "width", "contents.width":
		return &m.Width
	case "height", "contents.height":
		return &m.Height
	case "fps", "contents.fps":
		return &m.FPS
	case "fourcc", "contents.fourcc":
		return &m.Fourcc
	case "frames", "contents.frames":
		return &m.Frames
	case "created_at", "contents.created_at":
		return &m.CreatedAt
	case "updated_at", "contents.updated_at":
		return &m.UpdatedAt
	default:
		return nil
	}
}

func (m *Content) fieldPtrsByName(names []string) []interface{} {
	fields := []interface{}{}
	for _, n := range names {
		f := m.fieldPtrByName(n)
		fields = append(fields, f)
	}
	return fields
}

func (m *Content) isColumnName(name string) bool {
	for _, c := range m.columnNames() {
		if c == name {
			return true
		}
	}
	return false
}

func (m *Content) columnNames() []string {
	return []string{
		"id",
		"group_id",
		"name",
		"type",
		"path",
		"width",
		"height",
		"fps",
		"fourcc",
		"frames",
		"created_at",
		"updated_at",
	}
}
