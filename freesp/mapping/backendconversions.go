package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
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
		for _, path := range melem.PathList() {
			log.Printf("path %s\n", path)
			for _, mod := range freesp.ValidModes {
				log.Printf("mod %v\n", mod)
				mod := gr.CreatePathMode(path, mod)
				pos := melem.ModePosition(mod)
				if pos != empty {
					xmln.Entry = append(xmln.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
				}
			}
		}
		for i, p := range melem.(*mapelem).inports {
			xmlp := backend.XmlPortPosHintNew(melem.(*mapelem).node.InPorts()[i].Name())
			for _, path := range p.PathList() {
				for _, mod := range freesp.ValidModes {
					mod := gr.CreatePathMode(path, mod)
					pos := melem.ModePosition(mod)
					if pos != empty {
						xmlp.Entry = append(xmlp.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
					}
				}
			}
			xmln.InPorts = append(xmln.InPorts, *xmlp)
		}
		for i, p := range melem.(*mapelem).outports {
			xmlp := backend.XmlPortPosHintNew(melem.(*mapelem).node.OutPorts()[i].Name())
			for _, path := range p.PathList() {
				for _, mod := range freesp.ValidModes {
					mod := gr.CreatePathMode(path, mod)
					pos := melem.ModePosition(mod)
					if pos != empty {
						xmlp.Entry = append(xmlp.Entry, *backend.XmlModeHintEntryNew(string(mod), pos.X, pos.Y))
					}
				}
			}
			xmln.OutPorts = append(xmln.InPorts, *xmlp)
		}
		xmlm.MappedNodes = append(xmlm.MappedNodes, *xmln)
	}
	return

	/*
		for _, n := range m.Graph().ItsType().InputNodes() {
			xmlm.InputNodes = append(xmlm.InputNodes, *CreateXmlIONodePosHint(m, n, ""))
		}
		for _, n := range m.Graph().ItsType().OutputNodes() {
			xmlm.OutputNodes = append(xmlm.OutputNodes, *CreateXmlIONodePosHint(m, n, ""))
		}
		for _, n := range m.Graph().ItsType().ProcessingNodes() {
			hintlist := CreateXmlNodePosHint(m, n, "")
			for _, h := range hintlist {
				xmlm.MappedNodes = append(xmlm.MappedNodes, h)
			}
		}
		return
	*/
}

/*
func CreateXmlNodePosHint(m mp.MappingIf, nd bh.NodeIf, path string) (xmln []backend.XmlNodePosHint) {
	xmlnd := CreateXmlIONodePosHint(m, nd, path)
	xmlnd.Expanded = nd.Expanded()
	xmln = append(xmln, *xmlnd)
	nt := nd.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, n := range impl.Graph().ProcessingNodes() {
				var p, p2 string
				if len(path) == 0 {
					p = nd.Name()
				} else {
					p = fmt.Sprintf("%s/%s", path, nd.Name())
				}
				hintlist := CreateXmlNodePosHint(m, nd, p)
				for _, h := range hintlist {
					xmln = append(xmln, h)
				}
			}
			break
		}
	}
	return
}

func CreateXmlIONodePosHint(m mp.MappingIf, n bh.NodeIf, path string) (xmln *backend.XmlNodePosHint) {
	var p string
	if len(path) == 0 {
		p = n.Name()
	} else {
		p = fmt.Sprintf("%s/%s", path, n.Name())
	}
	xmln = backend.XmlNodePosHintNew(p)
	melem, ok := m.MappedElement(behaviour.NodeIdFromString(p, m.Graph().Filename()))
	if ok {
		empty := image.Point{}
		for _, mod := range freesp.ValidModes {
			xmlp := freesp.StringFromMode[mod]
			pos := melem.ModePosition(mod)
			if pos != empty {
				xmln.Entry = append(xmln.Entry, *backend.XmlModeHintEntryNew(xmlp, pos.X, pos.Y))
			}
		}
		for _, p := range n.InPorts() {
			xmlp := backend.XmlPortPosHintNew(p.Name())
			xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
			xmln.InPorts = append(xmln.InPorts, *xmlp)
		}
		for _, p := range n.OutPorts() {
			xmlp := backend.XmlPortPosHintNew(p.Name())
			xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
			xmln.OutPorts = append(xmln.OutPorts, *xmlp)
		}
	}
	return
}
*/

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
	} else { // descend node tree:
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
