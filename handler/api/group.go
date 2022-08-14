package api

type GroupView struct {
	Paging Paging `"json":"paging"`
}

type GroupViewReturn struct {
	Status
}

func newGroupView() Parameter {
	var gv GroupView
	return &gv
}

func (gv *GroupView) Processing() Return {
	var ret GroupViewReturn

	ret.success = true
	return &ret
}
