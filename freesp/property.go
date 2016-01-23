package freesp

import (
	tr "github.com/axel-freesp/sge/interface/tree"
)

type property uint

var _ tr.Property = property(0)

const (
	MayAddObject property = (1 << 0)
	MayEdit      property = (1 << 1)
	MayRemove    property = (1 << 2)
)

func (p property) IsReadOnly() bool {
	return p&(MayAddObject|MayEdit) == 0
}

func (p property) MayAddObject() bool {
	return p&MayAddObject != 0
}

func (p property) MayEdit() bool {
	return p&MayEdit != 0
}

func (p property) MayRemove() bool {
	return p&MayRemove != 0
}
