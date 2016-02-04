package platform

import (
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	pf "github.com/axel-freesp/sge/interface/platform"
)

// Hints file

func CreateXmlPlatformHint(p pf.PlatformIf) (xmlpl *backend.XmlPlatformHint) {
	xmlpl = backend.XmlPlatformHintNew(p.Filename())
	for _, a := range p.Arch() {
		xmlpl.Arch = append(xmlpl.Arch, *CreateXmlArchHint(a))
	}
	return
}

func CreateXmlArchHint(a pf.ArchIf) (xmla *backend.XmlArchPosHint) {
	xmla = backend.XmlArchPosHintNew(a.Name())
	xmla.Entry = freesp.CreateXmlModePosition(a).Entry
	for _, p := range a.ArchPorts() {
		xmla.ArchPorts = append(xmla.ArchPorts, *CreateXmlArchPortHint(p))
	}
	for _, p := range a.Processes() {
		xmla.Processes = append(xmla.Processes, *CreateXmlProcessHint(p))
	}
	return
}

func CreateXmlArchPortHint(p pf.ArchPortIf) (xmlp *backend.XmlArchPortPosHint) {
	xmlp = backend.XmlArchPortPosHintNew()
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlProcessHint(p pf.ProcessIf) (xmlp *backend.XmlProcessPosHint) {
	xmlp = backend.XmlProcessPosHintNew(p.Name())
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	for _, c := range p.InChannels() {
		xmlp.InChannels = append(xmlp.InChannels, *CreateXmlChannelHint(c))
	}
	for _, c := range p.OutChannels() {
		xmlp.OutChannels = append(xmlp.OutChannels, *CreateXmlChannelHint(c))
	}
	return
}

func CreateXmlChannelHint(c pf.ChannelIf) (xmlc *backend.XmlChannelPosHint) {
	xmlc = backend.XmlChannelPosHintNew()
	xmlc.Entry = freesp.CreateXmlModePosition(c).Entry
	return
}

// Platform model files

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
	//ret.Entry = freesp.CreateXmlModePosition(a).Entry
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
	//ret.Entry = freesp.CreateXmlModePosition(p).Entry
	return ret
}

func CreateXmlInChannel(ch pf.ChannelIf) *backend.XmlInChannel {
	ret := backend.XmlInChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	//ret.Entry = freesp.CreateXmlModePosition(ch).Entry
	//c := ch.(*channel)
	//ret.ArchPortHints.Entry = freesp.CreateXmlModePosition(c.archport).Entry
	return ret
}

func CreateXmlOutChannel(ch pf.ChannelIf) *backend.XmlOutChannel {
	ret := backend.XmlOutChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	//ret.Entry = freesp.CreateXmlModePosition(ch).Entry
	//c := ch.(*channel)
	//ret.ArchPortHints.Entry = freesp.CreateXmlModePosition(c.archport).Entry
	return ret
}
