package graph

import (
	"image"
	gr "github.com/axel-freesp/sge/interface/graph"
	bh "github.com/axel-freesp/sge/interface/behaviour"
)

type Port struct {
	SelectableBox
	userObj bh.PortIf
}

func PortNew(box image.Rectangle, userObj bh.PortIf) *Port {
	var config DrawConfig
	if userObj.Direction() == gr.InPort {
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

func (p Port) UserObj() bh.PortIf {
	return p.userObj
}

