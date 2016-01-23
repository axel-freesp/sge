package freesp

import (
	"io"
)

type FileDataIf interface {
	Filenamer
	io.Reader
	Write() (data []byte, err error)
	ReadFile(filepath string) error
	WriteFile(filepath string) error
}

type Filenamer interface {
	Filename() string
	SetFilename(string)
}
