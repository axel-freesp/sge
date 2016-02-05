package platform

import (
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	//"image"
)

type archPort struct {
	gr.ModePositionerObject
	channel pf.ChannelIf
}

var _ pf.ArchPortIf = (*archPort)(nil)

func archPortNew(ch pf.ChannelIf) *archPort {
	return &archPort{*gr.ModePositionerObjectNew(), ch}
}

func (p *archPort) Channel() pf.ChannelIf {
	return p.channel
}
