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
	var config DrawConfig
	if userObj.Direction() == interfaces.InPort {
		config = DrawConfig{ColorInit(ColorOption(InputPort)),
			ColorInit(ColorOption(HighlightInPort)),
			ColorInit(ColorOption(SelectInPort)),
			ColorInit(ColorOption(BoxFrame)),
			Color{},
			image.Point{}}
	} else {
		config = DrawConfig{ColorInit(ColorOption(OutputPort)),
			ColorInit(ColorOption(HighlightOutPort)),
			ColorInit(ColorOption(SelectOutPort)),
			ColorInit(ColorOption(BoxFrame)),
			Color{},
			image.Point{}}
	}
	return &Port{SelectableBoxInit(box, config), userObj}
}

var _ PortIf = (*Port)(nil)

func (p Port) UserObj() interfaces.PortObject {
	return p.userObj
}

