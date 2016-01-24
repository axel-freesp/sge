package platform

import (
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	pf "github.com/axel-freesp/sge/interface/platform"
)

func CreateXmlPlatform(p pf.PlatformIf) *backend.XmlPlatform {
	ret := backend.XmlPlatformNew()
	ret.PlatformId = p.PlatformId()
	for _, a := range p.Arch() {
		ret.Arch = append(ret.Arch, *CreateXmlArch(a))
	}
	return ret
}

func CreateXmlArch(a pf.ArchIf) *backend.XmlArch {
	ret := backend.XmlArchNew(a.Name())
	for _, t := range a.IOTypes() {
		ret.IOType = append(ret.IOType, *CreateXmlIOType(t))
	}
	for _, p := range a.Processes() {
		ret.Processes = append(ret.Processes, *CreateXmlProcess(p))
	}
	ret.Entry = freesp.CreateXmlModePosition(a).Entry
	return ret
}

func CreateXmlIOType(t pf.IOTypeIf) *backend.XmlIOType {
	return backend.XmlIOTypeNew(t.Name(), ioXmlModeMap[t.IOMode()])
}

func CreateXmlProcess(p pf.ProcessIf) *backend.XmlProcess {
	ret := backend.XmlProcessNew(p.Name())
	for _, c := range p.InChannels() {
		ret.InputChannels = append(ret.InputChannels, *CreateXmlInChannel(c))
	}
	for _, c := range p.OutChannels() {
		ret.OutputChannels = append(ret.OutputChannels, *CreateXmlOutChannel(c))
	}
	ret.Entry = freesp.CreateXmlModePosition(p).Entry
	return ret
}

func CreateXmlInChannel(ch pf.ChannelIf) *backend.XmlInChannel {
	ret := backend.XmlInChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = freesp.CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = freesp.CreateXmlModePosition(c.archport).Entry
	return ret
}

func CreateXmlOutChannel(ch pf.ChannelIf) *backend.XmlOutChannel {
	ret := backend.XmlOutChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = freesp.CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = freesp.CreateXmlModePosition(c.archport).Entry
	return ret
}
