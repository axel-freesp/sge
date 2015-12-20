package freesp

type property uint

var _ Property = property(0)

const (
	mayAddObject property = (1 << 0)
	mayEdit      property = (1 << 1)
	mayRemove    property = (1 << 2)
)

func (p property) IsReadOnly() bool {
	return p&(mayAddObject|mayEdit) == 0
}

func (p property) MayAddObject() bool {
	return p&mayAddObject != 0
}

func (p property) MayEdit() bool {
	return p&mayEdit != 0
}

func (p property) MayRemove() bool {
	return p&mayRemove != 0
}
