package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	mp "github.com/axel-freesp/sge/interface/mapping"
	"image"
)

func CreateXmlIOMap(node, process string, pos image.Point) *backend.XmlIOMap {
	return backend.XmlIOMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlNodeMap(node, process string, pos image.Point) *backend.XmlNodeMap {
	return backend.XmlNodeMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlMapping(m mp.MappingIf) (xmlm *backend.XmlMapping) {
	xmlm = backend.XmlMappingNew(m.Graph().Filename(), m.Platform().Filename())
	g := m.Graph().ItsType()
	for _, n := range g.InputNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.OutputNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.ProcessingNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.Mappings = append(xmlm.Mappings, *CreateXmlNodeMap(n.Name(), pname, melem.Position()))
		}
	}
	return
}
