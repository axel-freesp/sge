package freesp

import (
	"fmt"
	interfaces "github.com/axel-freesp/sge/interface"
	"image"
	"log"
)

type channel struct {
	direction interfaces.PortDirection
	iotype    IOType
	link      Channel
	process   Process
	linkText  string
	position  image.Point
	archPort  interfaces.ArchPortObject
}

var _ Channel = (*channel)(nil)
var _ interfaces.ChannelObject = (*channel)(nil)

func ChannelNew(dir interfaces.PortDirection, iotype IOType, process Process, linkText string, pos image.Point) *channel {
	return &channel{dir, iotype, nil, process, linkText, pos, nil}
}

func (c *channel) IOTypeObject() interfaces.IOTypeObject {
	return c.iotype.(*iotype)
}

func (c *channel) ProcessObject() interfaces.ProcessObject {
	return c.process.(*process)
}

func (c *channel) LinkObject() interfaces.ChannelObject {
	return c.link.(*channel)
}

func (c *channel) ArchPortObject() interfaces.ArchPortObject {
	return c.archPort
}

func (c *channel) Process() Process {
	return c.process
}

func (c *channel) IOType() IOType {
	return c.iotype
}

func (c *channel) SetIOType(newIOType IOType) {
	c.iotype = newIOType
	// update link
	if c.link != nil {
		c.link.(*channel).iotype = newIOType
	}
}

func (c *channel) Link() Channel {
	return c.link
}

/*
 *      Positioner API
 */

func (c *channel) Position() (p image.Point) {
	p = c.position
	return
}

func (c *channel) SetPosition(p image.Point) {
	c.position = p
}

/*
 *  Directioner API
 */

func (c *channel) Direction() interfaces.PortDirection {
	return c.direction
}

func (c *channel) SetDirection(newDir interfaces.PortDirection) {
	c.direction = newDir
}

/*
 *  Namer API
 */

func (c *channel) Name() string {
	if c.link != nil {
		return channelMakeName(c.iotype, c.link.Process())
	} else {
		return fmt.Sprintf("%s-%s", c.iotype.Name(), c.linkText)
	}
}

func channelMakeName(iotype IOType, link Process) string {
	return fmt.Sprintf("%s-%s/%s", iotype.Name(), link.Arch().Name(), link.Name())
}

func (c *channel) SetName(newName string) {
	log.Panicf("channel.SetName is forbidden.\n")
}

/*
 *  fmt.Stringer API
 */

func (c *channel) String() string {
	var dirtext string
	if c.Direction() == interfaces.InPort {
		dirtext = "in"
	} else {
		dirtext = "out"
	}
	return fmt.Sprintf("Channel(%s, %s, Link %s/%s, '%s')",
		dirtext, c.Name(), c.link.Process().Arch().Name(), c.link.Process().Name(), c.linkText)
}

/*
 *  TreeElement API
 */

func (c *channel) AddToTree(tree Tree, cursor Cursor) {
	var symbol Symbol
	if c.Direction() == interfaces.InPort {
		symbol = SymbolInChannel
	} else {
		symbol = SymbolOutChannel
	}
	err := tree.AddEntry(cursor, symbol, c.Name(), c, mayEdit|mayAddObject|mayRemove)
	if err != nil {
		log.Fatalf("channel.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (c *channel) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatalf("channel.AddNewObject error: nothing to add\n")
	return
}

func (c *channel) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatalf("channel.RemoveObject error: nothing to remove\n")
	return
}

/*
 *      channelList
 *
 */

type channelList struct {
	channels []Channel
	exports  []interfaces.ChannelObject
}

func channelListInit() channelList {
	return channelList{}
}

func (l *channelList) Append(ch *channel) {
	l.channels = append(l.channels, ch)
	l.exports = append(l.exports, ch)
}

func (l *channelList) Remove(ch Channel) {
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
		l.exports[i-1] = l.exports[i]
	}
	l.channels = l.channels[:len(l.channels)-1]
	l.exports = l.exports[:len(l.exports)-1]
}

func (l *channelList) Channels() []Channel {
	return l.channels
}

func (l *channelList) Exports() []interfaces.ChannelObject {
	return l.exports
}

func (l *channelList) Find(name string) (c Channel, ok bool) {
	for _, c = range l.channels {
		if c.Name() == name {
			ok = true
			return
		}
	}
	return
}
