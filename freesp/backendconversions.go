package freesp

import (
	"github.com/axel-freesp/sge/backend"
	gr "github.com/axel-freesp/sge/interface/graph"
	"image"
)

var ValidModes = []gr.PositionMode{
	gr.PositionModeNormal,
	gr.PositionModeMapping,
	gr.PositionModeExpanded,
}

var StringFromMode = map[gr.PositionMode]string{
	gr.PositionModeNormal:   "normal",
	gr.PositionModeMapping:  "mapping",
	gr.PositionModeExpanded: "expanded",
}

var ModeFromString = map[string]gr.PositionMode{
	"normal":   gr.PositionModeNormal,
	"mapping":  gr.PositionModeMapping,
	"expanded": gr.PositionModeExpanded,
}

func CreateXmlModePosition(x gr.ModePositioner) (h *backend.XmlModeHint) {
	h = backend.XmlModeHintNew()
	empty := image.Point{}
	for _, m := range ValidModes {
		pos := x.ModePosition(m)
		if pos != empty {
			h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(StringFromMode[m], pos.X, pos.Y))
		}
	}
	return
}
