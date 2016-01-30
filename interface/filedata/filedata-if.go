package filedata

import (
	"io"
)

type FileDataIf interface {
	Filenamer
	io.Reader
	Write() (data []byte, err error) // TODO: Make io.Writer
	ReadFile(filepath string) error
	WriteFile(filepath string) error
}

type Filenamer interface {
	Filename() string
	SetFilename(string)
	PathPrefix() string
	SetPathPrefix(string)
}
