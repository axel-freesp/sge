package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlNodeType struct {
	XMLName        xml.Name            `xml:"node-type"`
	TypeName       string              `xml:"name,attr"`
	InPort         []XmlInPort         `xml:"intype"`
	OutPort        []XmlOutPort        `xml:"outtype"`
	Implementation []XmlImplementation `xml:"implementation"`
}

func XmlNodeTypeNew(name string) *XmlNodeType {
	return &XmlNodeType{xml.Name{freespNamespace, "node-type"}, name, nil, nil, nil}
}

func (n *XmlNodeType) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, n)
	if err != nil {
		err = fmt.Errorf("XmlNodeType.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (n *XmlNodeType) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(n, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlNodeType.Write error: %v", err)
	}
	return
}
