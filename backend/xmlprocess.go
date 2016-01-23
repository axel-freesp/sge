package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlProcess struct {
	XMLName        xml.Name        `xml:"process"`
	Name           string          `xml:"name,attr"`
	InputChannels  []XmlInChannel  `xml:"input-channel"`
	OutputChannels []XmlOutChannel `xml:"output-channel"`
	XmlModeHint
}

func XmlProcessNew(name string) *XmlProcess {
	return &XmlProcess{xml.Name{freespNamespace, "process"}, name, nil, nil, XmlModeHint{}}
}

func (p *XmlProcess) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, p)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (p *XmlProcess) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(p, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}
