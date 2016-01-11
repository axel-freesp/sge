package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	interfaces "github.com/axel-freesp/sge/interface"
	"image"
	"log"
	"strings"
)

type platform struct {
	filename   string
	platformId string
	archlist   archList
	shape      image.Point
}

var _ Platform = (*platform)(nil)
var _ interfaces.PlatformObject = (*platform)(nil)

func PlatformNew(filename string) *platform {
	return &platform{filename, "", archListInit(), image.Point{}}
}

func (p *platform) PlatformId() string {
	return p.platformId
}

func (p *platform) SetPlatformId(newId string) {
	p.platformId = newId
}

func (p *platform) Arch() []Arch {
	return p.archlist.Archs()
}

func (p *platform) ArchObjects() []interfaces.ArchObject {
	return p.archlist.Exports()
}

func (p *platform) PlatformObject() interfaces.PlatformObject {
	return p
}

/*
 *  Shaper API
 */

func (p *platform) Shape() image.Point {
	return p.shape
}

func (p *platform) SetShape(shape image.Point) {
	p.shape = shape
}

/*
 *  Filenamer API
 */

func (p *platform) Filename() string {
	return p.filename
}

func (p *platform) SetFilename(newFilename string) {
	p.filename = newFilename
	// TODO
}

func (p *platform) Read(data []byte) (err error) {
	xmlp := backend.XmlPlatformNew()
	err = xmlp.Read(data)
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
	return
}

func (p *platform) createPlatformFromXml(xmlp *backend.XmlPlatform) (err error) {
	for _, xmla := range xmlp.Arch {
		a := ArchNew(xmla.Name, p)
		err = a.createArchFromXml(xmla)
		if err != nil {
			err = fmt.Errorf("platform.createPlatformFromXml error: %s", err)
			return
		}
		p.archlist.Append(a)
	}
	for _, a := range p.archlist.Archs() {
		for _, pr := range a.Processes() {
			for _, c := range pr.InChannels() {
				link := strings.Split(c.(*channel).linkText, "/")
				if len(link) != 2 {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s\n", c.(*channel).linkText)
					return
				}
				var ok bool
				var aa Arch
				var pp Process
				var cc Channel
				aa, ok = p.archlist.Find(link[0])
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no arch found)\n", c.(*channel).linkText)
					return
				}
				pp, ok = aa.(*arch).processes.Find(link[1])
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no process found)\n", c.(*channel).linkText)
					return
				}
				cName := channelMakeName(c.IOType(), pr)
				cc, ok = pp.(*process).outChannels.Find(cName)
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no channel found)\n", c.(*channel).linkText)
					return
				}
				c.(*channel).link = cc
			}
			for _, c := range pr.OutChannels() {
				link := strings.Split(c.(*channel).linkText, "/")
				if len(link) != 2 {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s\n", c.(*channel).linkText)
					return
				}
				var ok bool
				var aa Arch
				var pp Process
				var cc Channel
				aa, ok = p.archlist.Find(link[0])
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no arch found)\n", c.(*channel).linkText)
					return
				}
				pp, ok = aa.(*arch).processes.Find(link[1])
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no process found)\n", c.(*channel).linkText)
					return
				}
				cName := channelMakeName(c.IOType(), pr)
				cc, ok = pp.(*process).inChannels.Find(cName)
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no channel found)\n", c.(*channel).linkText)
					return
				}
				c.(*channel).link = cc
			}
		}
	}
	p.shape = image.Point{xmlp.Shape.W, xmlp.Shape.H}
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

/*
 *  TreeElement API
 */

func (p *platform) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolPlatform, p.Filename(), p, mayAddObject)
	if err != nil {
		log.Fatalf("platform.AddToTree error: AddEntry failed: %s", err)
	}
	for _, a := range p.Arch() {
		child := tree.Append(cursor)
		a.AddToTree(tree, child)
	}
}

func (p *platform) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("platform.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case Arch:
		a := obj.(Arch)
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

func (p *platform) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("platform.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Arch:
		a := obj.(Arch)
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
		removed = append(removed, IdWithObject{prefix, index, a})

	default:
		log.Fatalf("platform.RemoveObject error: invalid type %T\n", obj)
	}
	return
}
