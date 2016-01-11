package graph

import (
	"image"
	interfaces "github.com/axel-freesp/sge/interface"
)

type Port struct {
	SelectableBox
	userObj interfaces.PortObject
}

func PortNew(box image.Rectangle, userObj interfaces.PortObject) *Port {
	if userObj.Direction() == interfaces.InPort {
		return &Port{SelectableBoxInit(box,
			ColorInit(ColorOption(InputPort)),
			ColorInit(ColorOption(HighlightInPort)),
			ColorInit(ColorOption(SelectInPort)),
			ColorInit(ColorOption(BoxFrame)),
			image.Point{}),
			userObj}
	}
	return &Port{SelectableBoxInit(box,
			ColorInit(ColorOption(OutputPort)),
			ColorInit(ColorOption(HighlightOutPort)),
			ColorInit(ColorOption(SelectOutPort)),
			ColorInit(ColorOption(BoxFrame)),
			image.Point{}),
			userObj}
}

var _ PortIf = (*Port)(nil)

func (p Port) UserObj() interfaces.PortObject {
	return p.userObj
}

