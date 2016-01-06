package freesp

import (
	"fmt"
	"log"
	//"strings"
	"github.com/axel-freesp/sge/backend"
)

type process struct {
	name        string
	inChannels  channelList
	outChannels channelList
}

var _ Process = (*process)(nil)

func ProcessNew(name string) *process {
	return &process{name, channelListInit(), channelListInit()}
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
		c := ChannelNew(xmlc.Name, InPort, iot, p, xmlc.Source)
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
		c := ChannelNew(xmlc.Name, OutPort, iot, p, xmlc.Dest)
		p.outChannels.Append(c)
	}
	return
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
	return
}

func (p *process) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
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
