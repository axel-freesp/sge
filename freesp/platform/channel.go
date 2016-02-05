package platform

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	//"image"
	"log"
)

type channel struct {
	gr.ModePositionerObject
	direction gr.PortDirection
	iotype    pf.IOTypeIf
	link      pf.ChannelIf
	process   pf.ProcessIf
	linkText  string
	archport  pf.ArchPortIf
}

var _ pf.ChannelIf = (*channel)(nil)

func ChannelNew(dir gr.PortDirection, iotype pf.IOTypeIf, process pf.ProcessIf, linkText string) *channel {
	return &channel{*gr.ModePositionerObjectNew(), dir, iotype, nil, process, linkText, nil}
}

func createInChannelFromXml(xmlc backend.XmlInChannel, p pf.ProcessIf) (ch *channel, err error) {
	var iot pf.IOTypeIf
	iot, err = channelGetIOTypeFromArch(p.Arch(), xmlc.IOType)
	if err != nil {
		return
	}
	ch = ChannelNew(gr.InPort, iot, p, xmlc.Source)
	//ch.channelPositionsFromXml(xmlc.XmlChannel)
	ap := p.Arch().(*arch).AddArchPort(ch)
	//ap.archPortPositionsFromXml(xmlc.XmlChannel)
	ch.archport = ap
	return
}

func createOutChannelFromXml(xmlc backend.XmlOutChannel, p pf.ProcessIf) (ch *channel, err error) {
	var iot pf.IOTypeIf
	iot, err = channelGetIOTypeFromArch(p.Arch(), xmlc.IOType)
	if err != nil {
		return
	}
	ch = ChannelNew(gr.OutPort, iot, p, xmlc.Dest)
	//ch.channelPositionsFromXml(xmlc.XmlChannel)
	ap := p.Arch().(*arch).AddArchPort(ch)
	//ap.archPortPositionsFromXml(xmlc.XmlChannel)
	ch.archport = ap
	return
}

/*
func (ch *channel) channelPositionsFromXml(xmlc backend.XmlChannel) {
	for _, xmlh := range xmlc.Entry {
		mode, ok := freesp.ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("channelPositionsFromXml warning: hint mode %s not defined\n",
				xmlh.Mode)
			continue
		}
		ch.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
	}
	return
}

func (ap *archPort) archPortPositionsFromXml(xmlc backend.XmlChannel) {
	for _, xmlh := range xmlc.ArchPortHints.Entry {
		mode, ok := freesp.ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("archPortPositionsFromXml warning: hint mode %s not defined\n",
				xmlh.Mode)
			continue
		}
		ap.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
	}
	return
}
*/

func channelGetIOTypeFromArch(a pf.ArchIf, iotype string) (iot pf.IOTypeIf, err error) {
	var ok bool
	for _, iot = range a.IOTypes() {
		if iot.Name() == iotype {
			ok = true
			break
		}
	}
	if !ok {
		err = fmt.Errorf("createInChannelFromXml error: referenced ioType %s not found in arch %s.\n",
			iotype, a.Name())
	}
	return
}

func (c *channel) ArchPort() pf.ArchPortIf {
	return c.archport
}

func (c *channel) Process() pf.ProcessIf {
	return c.process
}

func (c *channel) IOType() pf.IOTypeIf {
	return c.iotype
}

func (c *channel) SetIOType(newIOType pf.IOTypeIf) {
	c.iotype = newIOType
	// update link
	if c.link != nil {
		c.link.(*channel).iotype = newIOType
	}
}

func (c *channel) Link() pf.ChannelIf {
	return c.link
}

func (c *channel) CreateXml() (buf []byte, err error) {
	if c.Direction() == gr.InPort {
		xmlc := CreateXmlInChannel(c)
		buf, err = xmlc.Write()
	} else {
		xmlc := CreateXmlOutChannel(c)
		buf, err = xmlc.Write()
	}
	return
}

//
//  Directioner API
//

func (c *channel) Direction() gr.PortDirection {
	return c.direction
}

func (c *channel) SetDirection(newDir gr.PortDirection) {
	c.direction = newDir
}

//
//  Namer API
//

func (c *channel) Name() string {
	if c.link != nil {
		return channelMakeName(c.iotype, c.link.Process())
	} else {
		return fmt.Sprintf("%s-%s", c.iotype.Name(), c.linkText)
	}
}

func channelMakeName(iotype pf.IOTypeIf, link pf.ProcessIf) string {
	return fmt.Sprintf("%s-%s/%s", iotype.Name(), link.Arch().Name(), link.Name())
}

func (c *channel) SetName(newName string) {
	log.Panicf("channel.SetName is forbidden.\n")
}

//
//  fmt.Stringer API
//

func (c *channel) String() string {
	var dirtext string
	if c.Direction() == gr.InPort {
		dirtext = "in"
	} else {
		dirtext = "out"
	}
	return fmt.Sprintf("Channel(%s, %s, Link %s/%s, '%s')",
		dirtext, c.Name(), c.link.Process().Arch().Name(), c.link.Process().Name(), c.linkText)
}

//
//  tr.TreeElement API
//

func (c *channel) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var symbol tr.Symbol
	if c.Direction() == gr.InPort {
		symbol = tr.SymbolInChannel
	} else {
		symbol = tr.SymbolOutChannel
	}
	prop := freesp.PropertyNew(true, true, true)
	err := tree.AddEntry(cursor, symbol, c.Name(), c, prop)
	if err != nil {
		log.Fatalf("channel.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (c *channel) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatalf("channel.AddNewObject error: nothing to add\n")
	return
}

func (c *channel) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatalf("channel.RemoveObject error: nothing to remove\n")
	return
}

func (c *channel) Identify(te tr.TreeElement) bool {
	switch te.(type) {
	case *channel:
		return te.(*channel).Name() == c.Name()
	}
	return false
}

//
//      channelList
//

type channelList struct {
	channels []pf.ChannelIf
}

func channelListInit() channelList {
	return channelList{}
}

func (l *channelList) Append(ch *channel) {
	l.channels = append(l.channels, ch)
}

func (l *channelList) Remove(ch pf.ChannelIf) {
	var i int
	for i = range l.channels {
		if ch == l.channels[i] {
			break
		}
	}
	if i >= len(l.channels) {
		for _, v := range l.channels {
			log.Printf("channelList.RemoveNodeType have Channel %v\n", v)
		}
		log.Fatalf("channelList.RemoveNodeType error: Channel %v not in this list\n", ch)
	}
	for i++; i < len(l.channels); i++ {
		l.channels[i-1] = l.channels[i]
	}
	l.channels = l.channels[:len(l.channels)-1]
}

func (l *channelList) Channels() []pf.ChannelIf {
	return l.channels
}

func (l *channelList) Find(name string) (c pf.ChannelIf, ok bool) {
	for _, c = range l.channels {
		if c.Name() == name {
			ok = true
			return
		}
	}
	return
}
