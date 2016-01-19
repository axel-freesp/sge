package freesp

import (
	interfaces "github.com/axel-freesp/sge/interface"
	"image"
)

type archPort struct {
	channel  interfaces.ChannelObject
	position map[interfaces.PositionMode]image.Point
}

var _ interfaces.ArchPortObject = (*archPort)(nil)

func archPortNew(ch interfaces.ChannelObject) *archPort {
	return &archPort{ch, make(map[interfaces.PositionMode]image.Point)}
}

func (p *archPort) Channel() interfaces.ChannelObject {
	return p.channel
}

/*
 *      ModePositioner API
 */

func (p archPort) ModePosition(mode interfaces.PositionMode) (pt image.Point) {
	pt = p.position[mode]
	return
}

func (p *archPort) SetModePosition(mode interfaces.PositionMode, pt image.Point) {
	p.position[mode] = pt
}
