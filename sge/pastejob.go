package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"log"
	"strings"
)

type PasteJob struct {
	context     string
	newElements []*NewElementJob
	children    []*PasteJob
}

func PasteJobNew() *PasteJob {
	return &PasteJob{"", nil, nil}
}

func (j *PasteJob) String() (text string) {
	text = fmt.Sprintf("PasteJob(")
	for _, e := range j.newElements {
		text = fmt.Sprintf("%s\n\t%s", text, e)
	}
	text = fmt.Sprintf("%s)", text)
	return
}

func ParseText(text string, fts *models.FilesTreeStore) (job *EditorJob, err error) {
	var parent freesp.TreeElement
	context := fts.GetCurrentId()
	if len(context) == 0 {
		err = fmt.Errorf("NewElementJob.ParseText TODO: Toplevel elements not implemented")
		return
	}
	parent, err = fts.GetObjectById(context)
	if err != nil {
		return
	}
	var ok bool
	var j *PasteJob
	switch parent.(type) {
	case freesp.SignalGraph:
		j, ok = parseNode(text, context, parent.(freesp.SignalGraph).ItsType())
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed Node\n")
			return
		}

	case freesp.SignalGraphType:
		j, ok = parseNode(text, context, parent.(freesp.SignalGraphType))
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed Node\n")
			return
		}

	case freesp.Node:
		j, ok = parseNode(text, getParentId(context), parent.(freesp.Node).Context())
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed Node\n")
			return
		}

	case freesp.NodeType:
		j, ok = parseNodeType(text, getParentId(context))
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed NodeType\n")
			return
		}

	case freesp.Port:

	case freesp.PortType:

	case freesp.Connection:

	case freesp.SignalType:
		j, ok = parseSignalType(text, getParentId(context))
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed SignalType\n")
			return
		}

	case freesp.Library:
		j, ok = parseNodeType(text, context)
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed NodeType\n")
			return
		}
		j, ok = parseSignalType(text, context)
		if ok {
			job = EditorJobNew(JobPaste, j)
			log.Printf("NewElementJob.ParseText: successfully parsed SignalType\n")
			return
		}

	case freesp.Implementation:

	default:
		err = fmt.Errorf("NewElementJob.ParseText error: can't insert to context %T", parent)
		return
	}

	err = fmt.Errorf("NewElementJob.ParseText error: function not implemented")
	return
}

func parseNode(text, context string, graph freesp.SignalGraphType) (job *PasteJob, ok bool) {
	xmln := backend.XmlNode{}
	xmlerr := xmln.Read([]byte(text))
	if xmlerr != nil {
		return
	}
	nodes := graph.Nodes()
	for true {
		valid := true
		for _, reg := range nodes {
			if reg.Name() == xmln.NName {
				valid = false
			}
		}
		if valid {
			break
		}
		xmln.NName = createNextNameCandidate(xmln.NName)
	}
	nodeTypes := freesp.GetRegisteredNodeTypes()
	validType := false
	for _, reg := range nodeTypes {
		if reg == xmln.NType {
			validType = true
		}
	}
	var njob *NewElementJob
	switch {
	case validType:
		njob = NewElementJobNew(context, eNode)
		njob.input[iNodeName] = xmln.NName
		njob.input[iNodeTypeSelect] = xmln.NType
	case len(xmln.InPort) == 0:
		if len(xmln.OutPort) == 0 {
			fmt.Printf("parseNode error: no ports.\n")
			return
		}
		njob = NewElementJobNew(context, eInputNode)
		njob.input[iInputNodeName] = xmln.NName
		njob.input[iInputTypeSelect] = xmln.OutPort[0].PType // TODO
	case len(xmln.OutPort) == 0:
		njob = NewElementJobNew(context, eOutputNode)
		njob.input[iOutputNodeName] = xmln.NName
		njob.input[iOutputTypeSelect] = xmln.InPort[0].PType // TODO
	default:
		if !validType {
			fmt.Printf("parseNode error: node type %s not registered.\n", xmln.NType)
			return
		}
		// TODO: Create node type??
		njob = NewElementJobNew(context, eNode)
		njob.input[iNodeName] = xmln.NName
		njob.input[iNodeTypeSelect] = xmln.NType
	}
	job = PasteJobNew()
	job.context = context
	job.newElements = append(job.newElements, njob)
	ok = true
	return
}

func parseNodeType(text, context string) (job *PasteJob, ok bool) {
	xmlnt := backend.XmlNodeType{}
	xmlerr := xmlnt.Read([]byte(text))
	if xmlerr != nil {
		return
	}
	nodeTypes := freesp.GetRegisteredNodeTypes()
	for true {
		valid := true
		for _, reg := range nodeTypes {
			if reg == xmlnt.TypeName {
				valid = false
			}
		}
		if valid {
			break
		}
		xmlnt.TypeName = createNextNameCandidate(xmlnt.TypeName)
	}
	job = PasteJobNew()
	job.context = context
	ntjob := NewElementJobNew(context, eNodeType)
	ntjob.input[iTypeName] = xmlnt.TypeName
	job.newElements = append(job.newElements, ntjob)
	for _, p := range xmlnt.InPort {
		pj := PasteJobNew()
		j := NewElementJobNew("", ePortType)
		j.input[iPortName] = p.PName
		j.input[iSignalTypeSelect] = p.PType
		j.input[iDirection] = direction2string[freesp.InPort]
		pj.newElements = append(pj.newElements, j)
		job.children = append(job.children, pj)
	}
	for _, p := range xmlnt.OutPort {
		pj := PasteJobNew()
		j := NewElementJobNew("", ePortType)
		j.input[iPortName] = p.PName
		j.input[iSignalTypeSelect] = p.PType
		j.input[iDirection] = direction2string[freesp.OutPort]
		pj.newElements = append(pj.newElements, j)
		job.children = append(job.children, pj)
	}
	for _, impl := range xmlnt.Implementation {
		pj := PasteJobNew()
		j := NewElementJobNew("", eImplementation)
		j.input[iImplName] = impl.Name
		if impl.SignalGraph == nil {
			j.input[iImplementationType] = implType2string[freesp.NodeTypeElement]
		} else {
			j.input[iImplementationType] = implType2string[freesp.NodeTypeGraph]
			// TODO: create graph objects
			// nodes: (no need for i/o nodes linked to ports...)
			for _, n := range impl.SignalGraph[0].ProcessingNodes {
				nj := NewElementJobNew("", eNode)
				nj.input[iNodeName] = n.NName
				validNodeType := false
				for _, nt := range nodeTypes {
					if nt == n.NType {
						validNodeType = true
					}
				}
				if !validNodeType {
					fmt.Printf("parseNodeType error: referenced node type %s not registered.\n", n.NType)
					return
				}
				nj.input[iNodeTypeSelect] = n.NType
				njob := PasteJobNew()
				njob.newElements = append(njob.newElements, nj)
				pj.children = append(pj.children, njob)
			}
			for _, e := range impl.SignalGraph[0].Connections {
				_ = e // TODO: how to navigate to port during apply???
			}
		}
		pj.newElements = append(pj.newElements, j)
		job.children = append(job.children, pj)
	}
	ok = true
	return
}

func parseSignalGraph(text, context string) (job *PasteJob, ok bool) {
	xmlsg := backend.XmlSignalGraph{}
	xmlerr := xmlsg.Read([]byte(text))
	if xmlerr != nil {
		return
	}
	nodeTypes := freesp.GetRegisteredNodeTypes()
	for _, l := range xmlsg.Libraries {
		_, err := global.GetLibrary(l.Name)
		if err != nil {
			fmt.Printf("parseSignalGraph error: referenced library %s not accessible.\n", l.Name)
			return
		}
	}
	var njob *NewElementJob
	njob = NewElementJobNew(context, eSignalGraph)
	job = PasteJobNew()
	job.context = context
	job.newElements = append(job.newElements, njob)
	for _, xmln := range xmlsg.ProcessingNodes {
		validNodeType := false
		for _, reg := range nodeTypes {
			if reg == xmln.NType {
				validNodeType = true
			}
		}
		if !validNodeType {
			fmt.Printf("parseSignalGraph error: node type %s not registered.\n", xmln.NType)
			return
		}
		njob = NewElementJobNew(context, eNode)
		njob.input[iNodeName] = xmln.NName
		njob.input[iNodeTypeSelect] = xmln.NType
		pj := PasteJobNew()
		pj.newElements = append(pj.newElements, njob)
		job.children = append(job.children, pj)
	}
	for _, xmln := range xmlsg.InputNodes {
		if len(xmln.OutPort) == 0 {
			fmt.Printf("parseSignalGraph error: input node has no outports.\n")
			return
		}
		njob = NewElementJobNew(context, eInputNode)
		njob.input[iInputNodeName] = xmln.NName
		njob.input[iInputTypeSelect] = xmln.OutPort[0].PType // TODO
		pj := PasteJobNew()
		pj.newElements = append(pj.newElements, njob)
		job.children = append(job.children, pj)
	}
	for _, xmln := range xmlsg.OutputNodes {
		if len(xmln.InPort) == 0 {
			fmt.Printf("parseSignalGraph error: output node has no inports.\n")
			return
		}
		njob = NewElementJobNew(context, eOutputNode)
		njob.input[iOutputNodeName] = xmln.NName
		njob.input[iOutputTypeSelect] = xmln.InPort[0].PType // TODO
		pj := PasteJobNew()
		pj.newElements = append(pj.newElements, njob)
		job.children = append(job.children, pj)
	}
	ok = true
	return
}

func parseSignalType(text, context string) (job *PasteJob, ok bool) {
	var scopeMap = map[string]string{
		"local":  "Local",
		"global": "Global",
		"":       "Global",
	}
	var modeMap = map[string]string{
		"sync":  "Isochronous",
		"async": "Asynchronous",
		"":      "Asynchronous",
	}
	xmlst := backend.XmlSignalType{}
	xmlerr := xmlst.Read([]byte(text))
	if xmlerr != nil {
		return
	}
	signalTypes := freesp.GetRegisteredSignalTypes()
	for true {
		valid := true
		for _, reg := range signalTypes {
			if reg == xmlst.Name {
				valid = false
			}
		}
		if valid {
			break
		}
		xmlst.Name = createNextNameCandidate(xmlst.Name)
	}
	job = PasteJobNew()
	job.context = context
	stjob := NewElementJobNew(context, eSignalType)
	stjob.input[iSignalTypeName] = xmlst.Name
	stjob.input[iCType] = xmlst.Ctype
	stjob.input[iChannelId] = xmlst.Msgid
	stjob.input[iScope] = scopeMap[xmlst.Scope]
	stjob.input[iSignalMode] = modeMap[xmlst.Mode]
	job.newElements = append(job.newElements, stjob)
	ok = true
	return
}

func createNextNameCandidate(text string) string {
	basename := baseName(text)
	suffix := text[len(basename):]
	suffixNum := 1
	if len(suffix) > 0 {
		fmt.Sscanf(suffix, "%d", &suffixNum)
	}
	return fmt.Sprintf("%s%d", basename, suffixNum+1)
}

func baseName(text string) string {
	suffix := text[len(text)-1:]
	if strings.ContainsAny(suffix, "0123456789") {
		return baseName(text[:len(text)-1])
	}
	return text
}
