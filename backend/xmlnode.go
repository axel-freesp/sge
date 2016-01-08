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
	Hint    XmlHint      `xml:"hint"`
}

type XmlHint struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

type XmlShape struct {
	W int `xml:"w,attr"`
	H int `xml:"h,attr"`
}

type XmlRectangle struct {
	XmlHint
	XmlShape
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

func (n *XmlNode) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, n)
	if err != nil {
		err = fmt.Errorf("XmlNode.Read error: %v", err)
	}
	return
}

func (n *XmlInputNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlInputNode.Write error: %v", err)
	}
	return
}

func XmlInputNodeNew(nName, nType string, x, y int) *XmlInputNode {
	return &XmlInputNode{xml.Name{freespNamespace, "input"}, "", XmlNode{nName, nType, nil, nil, XmlHint{x, y}}}
}

func (n *XmlOutputNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlOutputNode.Write error: %v", err)
	}
	return
}

func XmlOutputNodeNew(nName, nType string, x, y int) *XmlOutputNode {
	return &XmlOutputNode{xml.Name{freespNamespace, "output"}, "", XmlNode{nName, nType, nil, nil, XmlHint{x, y}}}
}

func (n *XmlProcessingNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlProcessingNode.Write error: %v", err)
	}
	return
}

func XmlProcessingNodeNew(nName, nType string, x, y int) *XmlProcessingNode {
	return &XmlProcessingNode{xml.Name{freespNamespace, "processing-node"}, XmlNode{nName, nType, nil, nil, XmlHint{x, y}}}
}
