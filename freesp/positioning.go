package freesp

import (
	"github.com/axel-freesp/sge/backend"
	gr "github.com/axel-freesp/sge/interface/graph"
	"image"
	"log"
)

func PathModePositionerApplyHints(pmp gr.PathModePositioner, xmln backend.XmlModeHint) {
	for _, m := range xmln.Entry {
		path, modestring := gr.SeparatePathMode(m.Mode)
		mode, ok := ModeFromString[string(modestring)]
		if !ok {
			log.Printf("PathModePositionerApplyHints Warning: hint mode %s not defined\n", m.Mode)
			continue
		}
		pos := image.Point{m.X, m.Y}
		pmp.SetPathModePosition(path, mode, pos)
	}
}

func ModePositionerApplyHints(mp gr.ModePositioner, xmln backend.XmlModeHint) {
	for _, m := range xmln.Entry {
		mode, ok := ModeFromString[m.Mode]
		if !ok {
			log.Printf("ModePositionerApplyHints Warning: hint mode %s not defined\n", m.Mode)
			continue
		}
		pos := image.Point{m.X, m.Y}
		mp.SetModePosition(mode, pos)
	}
}
