package backend

import (
	"encoding/xml"
    "os"
	"log"
)

type XmlBackend struct {
}

const envPathKey = "FREESP_PATH"
var xmlRoot string

func Init() {
    xmlRoot = os.Getenv(envPathKey)
    if len(xmlRoot) == 0 {
        log.Fatal("XmlBackendInit Error: missing environment variable FREESP_PATH")
    }
}

func XmlRoot() string {
    return xmlRoot
}

///////////////////////////////////////

type XmlLibraryRef struct {
    XMLName         xml.Name        `xml:"library"`
	Name string `xml:"ref,attr"`
}

