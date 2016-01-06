package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
	"strings"
)

type platform struct {
	filename   string
	platformId string
	archlist   archList
}

var _ Platform = (*platform)(nil)

func PlatformNew(filename string) *platform {
	return &platform{filename, "", archListInit()}
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

//func (p *platform) Read(data []byte) error

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
		a := ArchNew(xmla.Name)
		err = a.createArchFromXml(xmla)
		if err != nil {
			err = fmt.Errorf("platform.createPlatformFromXml error: %s", err)
			return
		}
		p.archlist.Append(a)
	}
	for _, a := range p.archlist.Archs() {
		for _, pr := range a.Processes() {
			ownLinkText := fmt.Sprintf("%s/%s", a.Name(), pr.Name())
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
				cc, ok = pp.(*process).outChannels.Find(ownLinkText)
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
				cc, ok = pp.(*process).inChannels.Find(ownLinkText)
				if !ok {
					err = fmt.Errorf("platform.createPlatformFromXml error (in): invalid source %s (no channel found)\n", c.(*channel).linkText)
					return
				}
				c.(*channel).link = cc
			}
		}
	}
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
	return
}

func (p *platform) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	return
}
