package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlNode struct {
	NName   string       `xml:"name,attr"`
	NType   string       `xml:"type,attr"`
	InPort  []XmlInPort  `xml:"intype"`
	OutPort []XmlOutPort `xml:"outtype"`
	XmlNodeHint
}

type XmlInputNode struct {
	XMLName xml.Name `xml:"input"`
	NPort   string   `xml:"port,attr"`
	XmlNode
}

type XmlOutputNode struct {
	XMLName xml.Name `xml:"output"`
	NPort   string   `xml:"port,attr"`
	XmlNode
}

type XmlProcessingNode struct {
	XMLName xml.Name `xml:"processing-node"`
	XmlNode
}

func (n *XmlNode) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, n)
	if err != nil {
		err = fmt.Errorf("XmlNode.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (n *XmlInputNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlInputNode.Write error: %v", err)
	}
	return
}

func XmlInputNodeNew(nName, nType string) *XmlInputNode {
	return &XmlInputNode{xml.Name{freespNamespace, "input"}, "", XmlNode{nName, nType, nil, nil, XmlNodeHint{XmlModeHint{}, false}}}
}

func (n *XmlOutputNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlOutputNode.Write error: %v", err)
	}
	return
}

func XmlOutputNodeNew(nName, nType string) *XmlOutputNode {
	return &XmlOutputNode{xml.Name{freespNamespace, "output"}, "", XmlNode{nName, nType, nil, nil, XmlNodeHint{XmlModeHint{}, false}}}
}

func (n *XmlProcessingNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlProcessingNode.Write error: %v", err)
	}
	return
}

func XmlProcessingNodeNew(nName, nType string, expanded bool) *XmlProcessingNode {
	return &XmlProcessingNode{xml.Name{freespNamespace, "processing-node"}, XmlNode{nName, nType, nil, nil, XmlNodeHint{XmlModeHint{}, expanded}}}
}
