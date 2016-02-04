package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlMappingHint struct {
	XMLName     xml.Name         `xml:"http://www.freesp.de/xml/freeSP hints"`
	Ref         string           `xml:"ref,attr"`
	MappedNodes []XmlNodePosHint `xml:"mapped-node"`
	Arch        []XmlArchPosHint `xml:"arch"`
}

func XmlMappingHintNew(ref string) *XmlMappingHint {
	return &XmlMappingHint{xml.Name{freespNamespace, "hints"}, ref, nil, nil}
}

func (h *XmlMappingHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlHint.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlMappingHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlHint.Write error: %v", err)
	}
	return
}
