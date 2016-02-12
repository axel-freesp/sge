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
	"strings"
)

type process struct {
	gr.ModePositionerObject
	name        string
	inChannels  channelList
	outChannels channelList
	arch        pf.ArchIf
}

var _ pf.ProcessIf = (*process)(nil)

func ProcessNew(name string, arch pf.ArchIf) *process {
	return &process{*gr.ModePositionerObjectNew(), name, channelListInit(), channelListInit(), arch}
}

func createProcessFromXml(xmlp backend.XmlProcess, a pf.ArchIf) (pr *process, err error) {
	pr = ProcessNew(xmlp.Name, a)
	for _, xmlc := range xmlp.InputChannels {
		var ch *channel
		ch, err = createInChannelFromXml(xmlc, pr)
		if err != nil {
			return
		}
		pr.inChannels.Append(ch)
	}
	for _, xmlc := range xmlp.OutputChannels {
		var ch *channel
		ch, err = createOutChannelFromXml(xmlc, pr)
		if err != nil {
			return
		}
		pr.outChannels.Append(ch)
	}
	/*
		for _, xmlh := range xmlp.Entry {
			mode, ok := freesp.ModeFromString[xmlh.Mode]
			if !ok {
				log.Printf("createProcessFromXml Warning: hint mode %s not defined\n",
					xmlh.Mode)
				continue
			}
			pr.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
		}
	*/
	return
}

func (p process) Arch() pf.ArchIf {
	return p.arch
}

func (p process) InChannels() []pf.ChannelIf {
	return p.inChannels.Channels()
}

func (p process) OutChannels() []pf.ChannelIf {
	return p.outChannels.Channels()
}

func (p process) CreateXml() (buf []byte, err error) {
	xmlp := CreateXmlProcess(&p)
	buf, err = xmlp.Write()
	return
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
 *  tr.TreeElementIf API
 */

func (p *process) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	prop := freesp.PropertyNew(true, true, true)
	err := tree.AddEntry(cursor, tr.SymbolProcess, p.Name(), p, prop)
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

func (p *process) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("process.AddNewObject error: %v nil object", p)
		return
	}
	switch obj.(type) {
	case pf.ChannelIf:
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
		var pp pf.ProcessIf
		pp, ok = aa.(*arch).processes.Find(link[1])
		if !ok {
			err = fmt.Errorf("process.AddNewObject error: %v invalid link text %s: no such process (abort)\n", p, c.linkText)
			return
		}
		var l *channelList
		var dd gr.PortDirection
		var ll *channelList
		var cPos, ccPos int
		if c.Direction() == gr.InPort {
			l = &p.inChannels
			ll = &pp.(*process).outChannels
			dd = gr.OutPort
			cPos = len(l.Channels())
			ccPos = tr.AppendCursor
		} else {
			l = &p.outChannels
			ll = &pp.(*process).inChannels
			dd = gr.InPort
			cPos = tr.AppendCursor
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
		if ppCursor.Position == tr.AppendCursor {
			ppCursor.Position = ccPos
		}
		newCursor = tree.Insert(ppCursor)
		cc.AddToTree(tree, newCursor)
		if cursor.Position == tr.AppendCursor {
			cursor.Position = cPos
		}
		newCursor = tree.Insert(cursor)
		c.AddToTree(tree, newCursor)
		if p.Arch().Name() != aa.Name() {
			p.Arch().(*arch).AddArchPort(c)
			aa.(*arch).AddArchPort(cc)
		}
		//log.Printf("process.AddNewObject: %v successfully added channel %v\n", p, c)

	default:
		log.Fatalf("process.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (p *process) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("process.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case pf.ChannelIf:
		c := obj.(pf.ChannelIf)
		cc := c.Link()
		pp := cc.Process()
		ppCursor := tree.Cursor(pp) // TODO. better search over platform...
		ccCursor := tree.CursorAt(ppCursor, cc)
		var l *channelList
		var ll *channelList
		if c.Direction() == gr.InPort {
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
		removed = append(removed, tr.IdWithObject{prefix, index, c})

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
	processs []pf.ProcessIf
}

func processListInit() processList {
	return processList{}
}

func (l *processList) Append(p *process) {
	l.processs = append(l.processs, p)
}

func (l *processList) Remove(p pf.ProcessIf) {
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

func (l *processList) Processes() []pf.ProcessIf {
	return l.processs
}

func (l *processList) Find(name string) (p pf.ProcessIf, ok bool) {
	for _, p = range l.processs {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}
