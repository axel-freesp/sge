package backend

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/sge/tool"
)

type XmlPlatform struct {
	XMLName    xml.Name  `xml:"http://www.freesp.de/xml/freeSP platform"`
	Version    string    `xml:"version,attr"`
	PlatformId string    `xml:"platform-id,attr"`
	Arch       []XmlArch `xml:"arch"`
	Shape      XmlShape  `xml:"hint"`
}

func XmlPlatformNew() *XmlPlatform {
	return &XmlPlatform{xml.Name{freespNamespace, "platform"}, "1.0", "", nil, XmlShape{}}
}

func (p *XmlPlatform) ReadFile(filepath string) (err error) {
	var data []byte
	data, err = tool.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("XmlPlatform.ReadFile error: Failed to read file %s", filepath)
		return
	}
	err = p.Read(data)
	if err != nil {
		err = fmt.Errorf("XmlPlatform.ReadFile error: %v", err)
	}
	return
}

func (p *XmlPlatform) WriteFile(filepath string) (err error) {
	var data []byte
	data, err = p.Write()
	if err != nil {
		return
	}
	buf := make([]byte, len(data)+len(xmlHeader))
	for i := 0; i < len(xmlHeader); i++ {
		buf[i] = xmlHeader[i]
	}
	for i := 0; i < len(data); i++ {
		buf[i+len(xmlHeader)] = data[i]
	}
	err = tool.WriteFile(filepath, buf)
	return
}

func (p *XmlPlatform) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, p)
	if err != nil {
		err = fmt.Errorf("XmlPlatform.Read error: %v", err)
	}
	return
}

func (p *XmlPlatform) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(p, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlPlatform.Write error: %v", err)
	}
	return
}
