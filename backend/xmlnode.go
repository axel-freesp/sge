package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlNode struct {
	NName   string    `xml:"name,attr"`
	NType   string    `xml:"type,attr"`
	InPort  []XmlInPort `xml:"intype"`
	OutPort []XmlOutPort `xml:"outtype"`
}

type XmlInputNode struct {
    XMLName         xml.Name        `xml:"input"`
    XmlNode
}

type XmlOutputNode struct {
    XMLName         xml.Name        `xml:"output"`
    XmlNode
}

type XmlProcessingNode struct {
    XMLName         xml.Name        `xml:"processing-node"`
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

func XmlInputNodeNew(nName, nType string) *XmlInputNode{
	return &XmlInputNode{xml.Name{"http://www.freesp.de/xml/freeSP", "input"}, XmlNode{nName, nType, nil, nil}}
}

func (n *XmlOutputNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
    if err != nil {
		err = fmt.Errorf("XmlOutputNode.Write error: %v", err)
    }
    return
}

func XmlOutputNodeNew(nName, nType string) *XmlOutputNode{
	return &XmlOutputNode{xml.Name{"http://www.freesp.de/xml/freeSP", "output"}, XmlNode{nName, nType, nil, nil}}
}

func (n *XmlProcessingNode) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
    if err != nil {
		err = fmt.Errorf("XmlProcessingNode.Write error: %v", err)
    }
    return
}

func XmlProcessingNodeNew(nName, nType string) *XmlProcessingNode{
	return &XmlProcessingNode{xml.Name{"http://www.freesp.de/xml/freeSP", "processing-node"}, XmlNode{nName, nType, nil, nil}}
}


