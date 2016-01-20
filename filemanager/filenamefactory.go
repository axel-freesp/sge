package filemanager

import (
	"fmt"
)

type filenameFactory struct {
	suffix string
	index  int
}

func FilenameFactoryInit(suffix string) filenameFactory {
	return filenameFactory{suffix, 0}
}

func (f *filenameFactory) NewFilename() (name string) {
	name = fmt.Sprintf("new-file-%d.%s", f.index, f.suffix)
	f.index++
	return
}
