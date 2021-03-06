package backend

import (
	"log"
	"os"
	"strings"
)

const envPathKey = "FREESP_PATH"
const envSearchPathKey = "FREESP_SEARCH_PATH"
const envSearchPathSep = ":"
const freespNamespace = "http://www.freesp.de/xml/freeSP"
const xmlHeader = `<?xml version="1.0" encoding="UTF-8"?>
`

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
