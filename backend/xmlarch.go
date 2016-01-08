package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlArch struct {
	XMLName   xml.Name     `xml:"arch"`
	Name      string       `xml:"name,attr"`
	IOType    []XmlIOType  `xml:"io-type"`
	Processes []XmlProcess `xml:"process"`
	Rect      XmlRectangle `xml:"hint"`
}

func XmlArchNew(name string) *XmlArch {
	return &XmlArch{xml.Name{freespNamespace, "arch"}, name, nil, nil, XmlRectangle{}}
}

func (a *XmlArch) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, a)
	if err != nil {
		err = fmt.Errorf("XmlArch.Read error: %v", err)
	}
	return
}

func (a *XmlArch) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(a, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlArch.Write error: %v", err)
	}
	return
}
