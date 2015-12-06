package backend

import (
	"encoding/xml"
	"log"
	"os"
	"strings"
)

const envPathKey = "FREESP_PATH"
const envSearchPathKey = "FREESP_SEARCH_PATH"
const envSearchPathSep = ":"
const freespNamespace = "http://www.freesp.de/xml/freeSP"

var xmlRoot string
var xmlSearchPaths []string

func Init() {
	xmlRoot = os.Getenv(envPathKey)
	if len(xmlRoot) == 0 {
		log.Fatal("XmlBackendInit Error: missing environment variable FREESP_PATH")
	}
	searchPath := os.Getenv(envSearchPathKey)
	xmlSearchPaths = strings.SplitN(searchPath, envSearchPathSep, -1)
	log.Println("backend.Init: search path:", XmlSearchPaths())
}

func XmlRoot() string {
	return xmlRoot
}

func XmlSearchPaths() []string {
	return xmlSearchPaths
}

///////////////////////////////////////

type XmlLibraryRef struct {
	XMLName xml.Name `xml:"library"`
	Name    string   `xml:"ref,attr"`
}

func XmlLibraryRefNew(filename string) *XmlLibraryRef {
	return &XmlLibraryRef{xml.Name{freespNamespace, "library"}, filename}
}
