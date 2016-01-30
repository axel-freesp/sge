package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlGraphHint struct {
	XMLName        xml.Name         `xml:"http://www.freesp.de/xml/freeSP hints"`
	Ref            string           `xml:"ref,attr"`
	InputNode      []XmlNodePosHint `xml:"input-node"`
	OutputNode     []XmlNodePosHint `xml:"output-node"`
	ProcessingNode []XmlNodePosHint `xml:"processing-node"`
}

func XmlGraphHintNew(ref string) *XmlGraphHint {
	return &XmlGraphHint{xml.Name{freespNamespace, "hints"}, ref, nil, nil, nil}
}

func (h *XmlGraphHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlHint.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlGraphHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlHint.Write error: %v", err)
	}
	return
}

type XmlNodePosHint struct {
	Name string `xml:"name,attr"`
	XmlModeHint
	Expanded bool             `xml:"expanded,attr"`
	InPorts  []XmlPortPosHint `xml:"in-port"`
	OutPorts []XmlPortPosHint `xml:"out-port"`
}

func XmlNodePosHintNew(name string) *XmlNodePosHint {
	return &XmlNodePosHint{name, XmlModeHint{}, false, nil, nil}
}

func (h *XmlNodePosHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlHint.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlNodePosHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlHint.Write error: %v", err)
	}
	return
}

type XmlPortPosHint struct {
	Name string `xml:"name,attr"`
	XmlModeHint
}

func XmlPortPosHintNew(name string) *XmlPortPosHint {
	return &XmlPortPosHint{name, XmlModeHint{}}
}

func (h *XmlPortPosHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlHint.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlPortPosHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlHint.Write error: %v", err)
	}
	return
}
