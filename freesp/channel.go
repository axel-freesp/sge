package freesp

import (
	//"fmt"
	"log"
)

type channel struct {
	name      string
	direction PortDirection
	iotype    IOType
	link      Channel
	process   Process
	linkText  string
}

var _ Channel = (*channel)(nil)

func ChannelNew(name string, dir PortDirection, iotype IOType, process Process, linkText string) *channel {
	return &channel{name, dir, iotype, nil, process, linkText}
}

func (c *channel) Process() Process {
	return c.process
}

func (c *channel) IOType() IOType {
	return c.iotype
}

func (c *channel) Link() Channel {
	return c.link
}

/*
 *  Directioner API
 */

func (c *channel) Direction() PortDirection {
	return c.direction
}

func (c *channel) SetDirection(newDir PortDirection) {
	c.direction = newDir
}

/*
 *  Namer API
 */

func (c *channel) Name() string {
	return c.name
}

func (c *channel) SetName(newName string) {
	c.name = newName
}

/*
 *  TreeElement API
 */

func (c *channel) AddToTree(tree Tree, cursor Cursor) {
	var symbol Symbol
	if c.Direction() == InPort {
		symbol = SymbolInChannel
	} else {
		symbol = SymbolOutChannel
	}
	err := tree.AddEntry(cursor, symbol, c.Name(), c, mayEdit|mayRemove)
	if err != nil {
		log.Fatalf("channel.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (c *channel) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	return
}

func (c *channel) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	return
}

/*
 *      channelList
 *
 */

type channelList struct {
	channels []Channel
}

func channelListInit() channelList {
	return channelList{nil}
}

func (l *channelList) Append(ch Channel) {
	l.channels = append(l.channels, ch)
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
	}
	l.channels = l.channels[:len(l.channels)-1]
}

func (l *channelList) Channels() []Channel {
	return l.channels
}

func (l *channelList) Find(linkText string) (c Channel, ok bool) {
	for _, c = range l.channels {
		if c.(*channel).linkText == linkText {
			ok = true
			return
		}
	}
	return
}
