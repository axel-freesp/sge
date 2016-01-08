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
	Rect           XmlRectangle    `xml:"hint"`
}

func XmlProcessNew(name string) *XmlProcess {
	return &XmlProcess{xml.Name{freespNamespace, "process"}, name, nil, nil, XmlRectangle{}}
}

func (p *XmlProcess) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, p)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	return
}

func (p *XmlProcess) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(p, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}
