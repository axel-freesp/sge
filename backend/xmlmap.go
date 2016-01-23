package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlMap struct {
	Name    string  `xml:"name,attr"`
	Process string  `xml:"process,attr"`
	Hint    XmlHint `xml:"hint"` // node position in map graphic
}

type XmlIOMap struct {
	XMLName xml.Name `xml:"map-ionode"`
	XmlMap
}

type XmlNodeMap struct {
	XMLName xml.Name `xml:"map-node"`
	XmlMap
}

func XmlIOMapNew(name, process string, x, y int) *XmlIOMap {
	return &XmlIOMap{xml.Name{freespNamespace, "map-ionode"}, XmlMap{name, process, XmlHint{x, y}}}
}

func XmlNodeMapNew(name, process string, x, y int) *XmlNodeMap {
	return &XmlNodeMap{xml.Name{freespNamespace, "map-node"}, XmlMap{name, process, XmlHint{x, y}}}
}

func (m *XmlIOMap) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, m)
	if err != nil {
		err = fmt.Errorf("XmlMap.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (m *XmlIOMap) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(m, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlMap.Write error: %v", err)
	}
	return
}

func (m *XmlNodeMap) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, m)
	if err != nil {
		err = fmt.Errorf("XmlMap.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (m *XmlNodeMap) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(m, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlMap.Write error: %v", err)
	}
	return
}
