package freesp

import (
	"github.com/axel-freesp/sge/backend"
	gr "github.com/axel-freesp/sge/interface/graph"
)

var validModes = []gr.PositionMode{gr.PositionModeNormal, gr.PositionModeMapping}

var StringFromMode = map[gr.PositionMode]string{
	gr.PositionModeNormal:  "normal",
	gr.PositionModeMapping: "mapping",
}

var ModeFromString = map[string]gr.PositionMode{
	"normal":  gr.PositionModeNormal,
	"mapping": gr.PositionModeMapping,
}

func CreateXmlModePosition(x gr.ModePositioner) (h *backend.XmlModeHint) {
	h = backend.XmlModeHintNew()
	for _, m := range validModes {
		pos := x.ModePosition(m)
		h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(StringFromMode[m], pos.X, pos.Y))
	}
	return
}
