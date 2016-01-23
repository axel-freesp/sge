package backend

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/sge/tool"
)

type XmlMapping struct {
	XMLName     xml.Name     `xml:"http://www.freesp.de/xml/freeSP mapping"`
	SignalGraph string       `xml:"graph,attr"`
	Platform    string       `xml:"platform,attr"`
	IOMappings  []XmlIOMap   `xml:"map-ionode"`
	Mappings    []XmlNodeMap `xml:"map-node"`
}

func XmlMappingNew(graph, platform string) *XmlMapping {
	return &XmlMapping{xml.Name{freespNamespace, "mapping"}, graph, platform, nil, nil}
}

func (m *XmlMapping) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, m)
	if err != nil {
		err = fmt.Errorf("XmlMap.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (m *XmlMapping) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(m, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlMap.Write error: %v", err)
	}
	return
}

func (m *XmlMapping) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("XmlMapping.ReadFile error: Failed to read file %s", filepath)
	}
	_, err = m.Read(data)
	if err != nil {
		return fmt.Errorf("XmlMapping.ReadFile error: %v", err)
	}
	return err
}

func (m *XmlMapping) WriteFile(filepath string) error {
	data, err := m.Write()
	if err != nil {
		return err
	}
	buf := make([]byte, len(data)+len(xmlHeader))
	for i := 0; i < len(xmlHeader); i++ {
		buf[i] = xmlHeader[i]
	}
	for i := 0; i < len(data); i++ {
		buf[i+len(xmlHeader)] = data[i]
	}
	return tool.WriteFile(filepath, buf)
}
