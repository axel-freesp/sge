package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/tool"
	"log"
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

func (f filenameFactory) HintFilename(filename string) (name string) {
	if tool.Suffix(filename) != f.suffix {
		log.Panicf("filenameFactory.HintFilename: invalid suffix\n")
		return
	}
	name = fmt.Sprintf("%s-%s.hints.xml", tool.Prefix(filename), f.suffix)
	return
}
