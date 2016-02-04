package platform

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
	"strings"
)

type platform struct {
	filename   string
	pathPrefix string
	platformId string
	archlist   archList
}

var _ pf.PlatformIf = (*platform)(nil)

func PlatformNew(filename string) *platform {
	return &platform{filename, "", "", archListInit()}
}

func PlatformApplyHints(p pf.PlatformIf, xmlhints *backend.XmlPlatformHint) (err error) {
	if p.Filename() != xmlhints.Ref {
		err = fmt.Errorf("PlatformApplyHints error: filename mismatch\n")
		return
	}
	for i, xmla := range xmlhints.Arch {
		a := p.Arch()[i]
		if a.Name() != xmla.Name {
			log.Printf("PlatformApplyHints error: arch name mismatch\n")
			continue
		}
		ArchApplyHints(a, &xmla)
	}
	return
}

func ArchApplyHints(a pf.ArchIf, xmla *backend.XmlArchPosHint) {
	freesp.ModePositionerApplyHints(a, xmla.XmlModeHint)
	for i, xmlp := range xmla.ArchPorts {
		p := a.ArchPorts()[i]
		freesp.ModePositionerApplyHints(p, xmlp.XmlModeHint)
	}
	for i, xmlp := range xmla.Processes {
		p := a.Processes()[i]
		ProcessApplyHints(p, &xmlp)
	}
}

func ProcessApplyHints(p pf.ProcessIf, xmlp *backend.XmlProcessPosHint) {
	freesp.ModePositionerApplyHints(p, xmlp.XmlModeHint)
	for i, xmlc := range xmlp.InChannels {
		c := p.InChannels()[i]
		freesp.ModePositionerApplyHints(c, xmlc.XmlModeHint)
	}
	for i, xmlc := range xmlp.OutChannels {
		c := p.OutChannels()[i]
		freesp.ModePositionerApplyHints(c, xmlc.XmlModeHint)
	}
}

func (p *platform) PlatformId() string {
	return p.platformId
}

func (p *platform) SetPlatformId(newId string) {
	p.platformId = newId
}

func (p *platform) Arch() []pf.ArchIf {
	return p.archlist.Archs()
}

func (p *platform) ProcessByName(name string) (pr pf.ProcessIf, ok bool) {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return
	}
	var a pf.ArchIf
	for _, a = range p.Arch() {
		if a.Name() == parts[0] {
			break
		}
	}
	if a == nil {
		return
	}
	for _, pr = range a.Processes() {
		if pr.Name() == parts[1] {
			ok = true
			return
		}
	}
	return
}

func (p *platform) CreateXml() (buf []byte, err error) {
	xmlp := CreateXmlPlatform(p)
	buf, err = xmlp.Write()
	return
}

//
//  Filenamer API
//

func (p *platform) Filename() string {
	return p.filename
}

func (p *platform) SetFilename(newFilename string) {
	p.filename = newFilename
	// TODO
}

func (p platform) PathPrefix() string {
	return p.pathPrefix
}

func (p *platform) SetPathPrefix(newP string) {
	p.pathPrefix = newP
}

func (p *platform) Read(data []byte) (cnt int, err error) {
	xmlp := backend.XmlPlatformNew()
	cnt, err = xmlp.Read(data)
	if err != nil {
		err = fmt.Errorf("platform.Read: %v", err)
		return
	}
	err = p.createPlatformFromXml(xmlp)
	return
}

func (p *platform) ReadFile(filepath string) (err error) {
	xmlp := backend.XmlPlatformNew()
	err = xmlp.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("platform.Read: %v", err)
		return
	}
	err = p.createPlatformFromXml(xmlp)
	if err != nil {
		log.Printf("platform.ReadFile error: %s\n", err)
	}
	return
}

func (p *platform) FindLinkedChannel(c1 pf.ChannelIf) (c2 pf.ChannelIf, err error) {
	linktext := c1.(*channel).linkText
	p1 := c1.Process()
	link := strings.Split(linktext, "/")
	if len(link) != 2 {
		err = fmt.Errorf("platform.FindLinkedChannel error: invalid reference %s\n", linktext)
		return
	}
	var ok bool
	var a2 pf.ArchIf
	var p2 pf.ProcessIf
	a2, ok = p.archlist.Find(link[0])
	if !ok {
		err = fmt.Errorf("platform.FindLinkedChannel error: invalid reference %s (no arch found)\n", linktext)
		return
	}
	p2, ok = a2.(*arch).processes.Find(link[1])
	if !ok {
		err = fmt.Errorf("platform.FindLinkedChannel error: invalid reference %s (no process found)\n", linktext)
		return
	}
	cName := channelMakeName(c1.IOType(), p1)
	if c1.Direction() == gr.InPort {
		c2, ok = p2.(*process).outChannels.Find(cName)
	} else {
		c2, ok = p2.(*process).inChannels.Find(cName)
	}
	if !ok {
		err = fmt.Errorf("platform.FindLinkedChannel error: invalid reference %s (no channel found)\n", linktext)
		return
	}
	return
}

func (p *platform) createPlatformFromXml(xmlp *backend.XmlPlatform) (err error) {
	for _, xmla := range xmlp.Arch {
		var a *arch
		a, err = createArchFromXml(xmla, p)
		if err != nil {
			return
		}
		p.archlist.Append(a)
	}
	for _, a := range p.archlist.Archs() {
		for _, pr := range a.Processes() {
			for _, c := range pr.InChannels() {
				c.(*channel).link, err = p.FindLinkedChannel(c)
				if err != nil {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): %s", err)
					return
				}
			}
			for _, c := range pr.OutChannels() {
				c.(*channel).link, err = p.FindLinkedChannel(c)
				if err != nil {
					err = fmt.Errorf("platform.createPlatformFromXml error (out): %s", err)
					return
				}
			}
		}
	}
	log.Printf("platform.createPlatformFromXml: all channels linked\n")
	return
}

func (p *platform) Write() (data []byte, err error) {
	xmlp := CreateXmlPlatform(p)
	data, err = xmlp.Write()
	return
}

func (p *platform) WriteFile(filepath string) (err error) {
	xmlp := CreateXmlPlatform(p)
	err = xmlp.WriteFile(filepath)
	return
}

func (p *platform) RemoveFromTree(tree tr.TreeIf) {
	tree.Remove(tree.Cursor(p))
}

//
//  tr.TreeElement API
//

func (p *platform) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	prop := freesp.PropertyNew(true, false, false)
	err := tree.AddEntry(cursor, tr.SymbolPlatform, p.Filename(), p, prop)
	if err != nil {
		log.Fatalf("platform.AddToTree error: AddEntry failed: %s", err)
	}
	for _, a := range p.Arch() {
		child := tree.Append(cursor)
		a.AddToTree(tree, child)
	}
}

func (p *platform) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("platform.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case pf.ArchIf:
		a := obj.(pf.ArchIf)
		_, ok := p.archlist.Find(a.Name())
		if ok {
			err = fmt.Errorf("platform.AddNewObject warning: duplicate arch name %s (abort)\n", a.Name())
			return
		}
		p.archlist.Append(a.(*arch))
		newCursor = tree.Insert(cursor)
		a.AddToTree(tree, newCursor)

	default:
		log.Fatalf("platform.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (p *platform) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("platform.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case pf.ArchIf:
		a := obj.(pf.ArchIf)
		_, ok := p.archlist.Find(a.Name())
		if ok {
			p.archlist.Remove(a)
		} else {
			log.Printf("platform.RemoveObject warning: arch name %s not found\n", a.Name())
		}
		for _, pr := range a.Processes() {
			prCursor := tree.CursorAt(cursor, pr)
			del := a.RemoveObject(tree, prCursor)
			for _, d := range del {
				removed = append(removed, d)
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, a})

	default:
		log.Fatalf("platform.RemoveObject error: invalid type %T\n", obj)
	}
	return
}
