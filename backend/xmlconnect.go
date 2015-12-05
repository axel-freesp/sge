package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlConnect struct {
    XMLName         xml.Name        `xml:"connect"`
	From     string `xml:"from,attr"`
	To       string `xml:"to,attr"`
	FromPort string `xml:"from-port,attr"`
	ToPort   string `xml:"to-port,attr"`
}

func (c *XmlConnect) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, c)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	return
}

func (c *XmlConnect) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(c, "", "   ")
    if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
    }
    return
}

