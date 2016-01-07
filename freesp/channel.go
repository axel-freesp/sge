package freesp

import (
	"fmt"
	"log"
)

type channel struct {
	direction PortDirection
	iotype    IOType
	link      Channel
	process   Process
	linkText  string
}

var _ Channel = (*channel)(nil)

func ChannelNew(dir PortDirection, iotype IOType, process Process, linkText string) *channel {
	return &channel{dir, iotype, nil, process, linkText}
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
	if c.Direction() == InPort {
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
	if c.Direction() == InPort {
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

func (l *channelList) Find(name string) (c Channel, ok bool) {
	for _, c = range l.channels {
		if c.Name() == name {
			ok = true
			return
		}
	}
	return
}
