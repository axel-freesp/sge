package freesp

import (
	"github.com/axel-freesp/sge/backend"
	gr "github.com/axel-freesp/sge/interface/graph"
	"image"
	//"fmt"
)

func CreateXmlModePosition(x gr.ModePositioner) (h *backend.XmlModeHint) {
	h = backend.XmlModeHintNew()
	empty := image.Point{}
	for _, m := range gr.ValidModes {
		pos := x.ModePosition(m)
		if pos != empty {
			h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(gr.StringFromMode[m], pos.X, pos.Y))
		}
	}
	return
}

/*
func AddXmlPathModePosition(h *backend.XmlModeHint, path string, x gr.PathModePositioner) {
	empty := image.Point{}
	for _, m := range ValidModes {
		pos := x.PathModePosition(path, m)
		if pos != empty {
			p := fmt.Sprintf("%s/%s", path, StringFromMode[m])
			h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(p, pos.X, pos.Y))
		}
	}
	return
}
*/
