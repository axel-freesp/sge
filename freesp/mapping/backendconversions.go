package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	//"github.com/axel-freesp/sge/freesp"
	//	"github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	"image"
	"log"
)

func CreateXmlMappingHint(m mp.MappingIf) (xmlm *backend.XmlMappingHint) {
	xmlm = backend.XmlMappingHintNew(m.Filename())
	for _, nId := range m.MappedIds() {
		melem, ok := m.MappedElement(nId)
		if !ok {
			log.Fatal("CreateXmlMappingHint internal error: inconsistent maplist.\n")
		}
		log.Printf("CreateXmlMappingHint(%s): nId=%s, melem.mode=%v, pos=%v\n", melem.NodeId(), nId.String(), melem.ActiveMode(), melem.Position())
		xmln := backend.XmlNodePosHintNew(nId.String())
		empty := image.Point{}
		for _, mod := range gr.ValidModes {
			pos := melem.ModePosition(mod)
			log.Printf("mod %v, pos=%v\n", mod, pos)
			if pos != empty {
				xmln.Entry = append(xmln.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
			}
		}
		for i, p := range melem.(*mapelem).inports {
			xmlp := backend.XmlPortPosHintNew(melem.(*mapelem).node.InPorts()[i].Name())
			log.Printf("CreateXmlMappingHint(inports):    xmlp=%s\n", xmlp.Name)
			for _, mod := range gr.ValidModes {
				pos := p.ModePosition(mod)
				log.Printf("   mod %v, pos=%v\n", mod, pos)
				if pos != empty {
					xmlp.Entry = append(xmlp.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
				}
			}
			xmln.InPorts = append(xmln.InPorts, *xmlp)
		}
		for i, p := range melem.(*mapelem).outports {
			xmlp := backend.XmlPortPosHintNew(melem.(*mapelem).node.OutPorts()[i].Name())
			log.Printf("CreateXmlMappingHint(outports):    xmlp=%s\n", xmlp.Name)
			for _, mod := range gr.ValidModes {
				pos := p.ModePosition(mod)
				log.Printf("   mod %v, pos=%v\n", mod, pos)
				if pos != empty {
					xmlp.Entry = append(xmlp.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
				}
			}
			xmln.OutPorts = append(xmln.InPorts, *xmlp)
		}
		xmlm.MappedNodes = append(xmlm.MappedNodes, *xmln)
	}
	return
}

func CreateXmlIOMap(node, process string) *backend.XmlIOMap {
	return backend.XmlIOMapNew(node, process)
}

func CreateXmlNodeMap(node, process string) *backend.XmlNodeMap {
	return backend.XmlNodeMapNew(node, process)
}

func CreateXmlNodeMapList(m mp.MappingIf, n bh.NodeIf, path string) (xmln []backend.XmlNodeMap) {
	p, ok := m.Mapped(n.Name())
	if ok { // entire node is mapped to p:
		pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
		xmln = append(xmln, *CreateXmlNodeMap(path, pname))
	}
	nt := n.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, nn := range impl.Graph().ProcessingNodes() {
				xmlnn := CreateXmlNodeMapList(m, nn, fmt.Sprintf("%s/%s", path, nn.Name()))
				for _, x := range xmlnn {
					xmln = append(xmln, x)
				}
			}
		}
	}
	return
}

func CreateXmlMapping(m mp.MappingIf) (xmlm *backend.XmlMapping) {
	xmlm = backend.XmlMappingNew(m.Graph().Filename(), m.Platform().Filename())
	g := m.Graph().ItsType()
	for _, n := range g.InputNodes() {
		p, ok := m.Mapped(n.Name())
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname))
		}
	}
	for _, n := range g.OutputNodes() {
		p, ok := m.Mapped(n.Name())
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname))
		}
	}
	for _, n := range g.ProcessingNodes() {
		xmln := CreateXmlNodeMapList(m, n, n.Name())
		for _, x := range xmln {
			xmlm.Mappings = append(xmlm.Mappings, x)
		}
	}
	return
}
