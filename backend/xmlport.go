package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlPort struct {
	PName string `xml:"port,attr"`
	PType string `xml:"type,attr"`
	XmlModeHint
}

type XmlInPort struct {
	XMLName xml.Name `xml:"intype"`
	XmlPort
}

func XmlInPortNew(pName, pType string) (xmlp *XmlInPort) {
	xmlp = &XmlInPort{xml.Name{freespNamespace, "intype"},
		XmlPort{pName, pType, XmlModeHint{}}}
	return
}

type XmlOutPort struct {
	XMLName xml.Name `xml:"outtype"`
	XmlPort
}

func XmlOutPortNew(pName, pType string) (xmlp *XmlOutPort) {
	xmlp = &XmlOutPort{xml.Name{freespNamespace, "outtype"},
		XmlPort{pName, pType, XmlModeHint{}}}
	return
}

func (p *XmlPort) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, p)
	if err != nil {
		err = fmt.Errorf("XmlPort.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (p *XmlInPort) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(p, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlInPort.Write error: %v", err)
	}
	return
}

func (p *XmlOutPort) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(p, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlOutPort.Write error: %v", err)
	}
	return
}
