package freesp

import (
	"github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	"image"
)

type archPort struct {
	channel  pf.ChannelIf
	position map[graph.PositionMode]image.Point
}

var _ pf.ArchPortIf = (*archPort)(nil)

func archPortNew(ch pf.ChannelIf) *archPort {
	return &archPort{ch, make(map[graph.PositionMode]image.Point)}
}

func (p *archPort) Channel() pf.ChannelIf {
	return p.channel
}

/*
 *      ModePositioner API
 */

func (p archPort) ModePosition(mode graph.PositionMode) (pt image.Point) {
	pt = p.position[mode]
	return
}

func (p *archPort) SetModePosition(mode graph.PositionMode, pt image.Point) {
	p.position[mode] = pt
}
