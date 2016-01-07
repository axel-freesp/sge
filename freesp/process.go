package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
	"strings"
)

type process struct {
	name        string
	inChannels  channelList
	outChannels channelList
	arch        Arch
}

var _ Process = (*process)(nil)

func ProcessNew(name string, arch Arch) *process {
	return &process{name, channelListInit(), channelListInit(), arch}
}

func (p *process) createProcessFromXml(xmlp backend.XmlProcess, ioTypes []IOType) (err error) {
	p.name = xmlp.Name
	for _, xmlc := range xmlp.InputChannels {
		var iot IOType
		ok := false
		for _, iot = range ioTypes {
			if iot.Name() == xmlc.IOType {
				ok = true
				break
			}
		}
		if !ok {
			err = fmt.Errorf("process.createProcessFromXml error (in): referenced ioType %s not found.\n", xmlc.IOType)
			return
		}
		c := ChannelNew(InPort, iot, p, xmlc.Source)
		p.inChannels.Append(c)
	}
	for _, xmlc := range xmlp.OutputChannels {
		var iot IOType
		ok := false
		for _, iot = range ioTypes {
			if iot.Name() == xmlc.IOType {
				ok = true
				break
			}
		}
		if !ok {
			err = fmt.Errorf("process.createProcessFromXml error (out): referenced ioType %s not found.\n", xmlc.IOType)
			return
		}
		c := ChannelNew(OutPort, iot, p, xmlc.Dest)
		p.outChannels.Append(c)
	}
	return
}

func (p *process) Arch() Arch {
	return p.arch
}

func (p *process) InChannels() []Channel {
	return p.inChannels.Channels()
}

func (p *process) OutChannels() []Channel {
	return p.outChannels.Channels()
}

/*
 *  Namer API
 */

func (p *process) Name() string {
	return p.name
}

func (p *process) SetName(newName string) {
	p.name = newName
}

/*
 *  fmt.Stringer API
 */

func (p *process) String() string {
	return fmt.Sprintf("Process(%s/%s)", p.arch.Name(), p.name)
}

/*
 *  TreeElement API
 */

func (p *process) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolProcess, p.Name(), p, mayAddObject|mayEdit|mayRemove)
	if err != nil {
		log.Fatalf("process.AddToTree error: AddEntry failed: %s\n", err)
	}
	for _, c := range p.InChannels() {
		child := tree.Append(cursor)
		c.AddToTree(tree, child)
	}
	for _, c := range p.OutChannels() {
		child := tree.Append(cursor)
		c.AddToTree(tree, child)
	}
}

func (p *process) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("process.AddNewObject error: %v nil object", p)
		return
	}
	switch obj.(type) {
	case Channel:
		c := obj.(*channel)
		c.process = p
		cLinkText := c.linkText
		link := strings.Split(cLinkText, "/")
		if len(link) != 2 {
			err = fmt.Errorf("process.AddNewObject error: %v invalid link text %s (abort)\n", p, c.linkText)
			return
		}
		archlist := p.Arch().Platform().(*platform).archlist
		aa, ok := archlist.Find(link[0])
		if !ok {
			err = fmt.Errorf("process.AddNewObject error: %v invalid link text %s: no such arch (abort)\n", p, c.linkText)
			return
		}
		var pp Process
		pp, ok = aa.(*arch).processes.Find(link[1])
		if !ok {
			err = fmt.Errorf("process.AddNewObject error: %v invalid link text %s: no such process (abort)\n", p, c.linkText)
			return
		}
		var l *channelList
		var dd PortDirection
		var ll *channelList
		var cPos, ccPos int
		if c.Direction() == InPort {
			l = &p.inChannels
			ll = &pp.(*process).outChannels
			dd = OutPort
			cPos = len(l.Channels())
			ccPos = AppendCursor
		} else {
			l = &p.outChannels
			ll = &pp.(*process).inChannels
			dd = InPort
			cPos = AppendCursor
			ccPos = len(ll.Channels())
		}
		cName := channelMakeName(c.iotype, pp)
		_, ok = l.Find(cName)
		if ok {
			err = fmt.Errorf("process.AddNewObject warning: %v duplicate %v channel name %s/%s (abort)\n",
				p, c.Direction(), c.linkText, c.iotype.Name())
			return
		}
		ccLinkText := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
		ccName := channelMakeName(c.iotype, p)
		_, ok = l.Find(ccName)
		_, ok = ll.Find(cName)
		if ok {
			err = fmt.Errorf("process.AddNewObject warning: %v duplicate %v channel name %s/%s on other side (abort)\n",
				pp, dd, ccLinkText, c.iotype.Name())
			return
		}
		cc := ChannelNew(dd, c.iotype, pp, ccLinkText)
		cc.link = c
		c.link = cc
		l.Append(c)
		ll.Append(cc)
		ppCursor := tree.Cursor(pp)
		if ppCursor.Position == AppendCursor {
			ppCursor.Position = ccPos
		}
		newCursor = tree.Insert(ppCursor)
		cc.AddToTree(tree, newCursor)
		if cursor.Position == AppendCursor {
			cursor.Position = cPos
		}
		newCursor = tree.Insert(cursor)
		c.AddToTree(tree, newCursor)
		//log.Printf("process.AddNewObject: %v successfully added channel %v\n", p, c)

	default:
		log.Fatalf("process.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (p *process) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("process.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Channel:
		c := obj.(Channel)
		cc := c.Link()
		pp := cc.Process()
		ppCursor := tree.Cursor(pp) // TODO. better search over platform...
		ccCursor := tree.CursorAt(ppCursor, cc)
		var l *channelList
		var ll *channelList
		if c.Direction() == InPort {
			l = &p.inChannels
			ll = &pp.(*process).outChannels
		} else {
			l = &p.outChannels
			ll = &pp.(*process).inChannels
		}
		l.Remove(c)
		ll.Remove(cc)
		tree.Remove(ccCursor)
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, c})

	default:
		log.Fatalf("Port.RemoveObject error: invalid type %T: %v\n", obj, obj)
	}
	return
}

/*
 *      processList
 *
 */

type processList struct {
	processs []Process
}

func processListInit() processList {
	return processList{nil}
}

func (l *processList) Append(p Process) {
	l.processs = append(l.processs, p)
}

func (l *processList) Remove(p Process) {
	var i int
	for i = range l.processs {
		if p == l.processs[i] {
			break
		}
	}
	if i >= len(l.processs) {
		for _, v := range l.processs {
			log.Printf("processList.RemoveNodeType have Process %v\n", v)
		}
		log.Fatalf("processList.RemoveNodeType error: Process %v not in this list\n", p)
	}
	for i++; i < len(l.processs); i++ {
		l.processs[i-1] = l.processs[i]
	}
	l.processs = l.processs[:len(l.processs)-1]
}

func (l *processList) Processes() []Process {
	return l.processs
}

func (l *processList) Find(name string) (p Process, ok bool) {
	for _, p = range l.processs {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}
