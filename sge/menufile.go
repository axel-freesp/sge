package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	fd "github.com/axel-freesp/sge/interface/filedata"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
)

type DirectoryMgrIf interface {
	// DirList(): To iterate through when resolving references
	DirList() []string
	// Current(): Set for open/save dialogs
	Current(ft FileType) string
	// SetCurrent(): The actual directory chosen by dialog
	SetCurrent(ft FileType, dirname string)
	// What is stored as obj.Filename()
	FilenameToShow(string) (prefix, filename string)
}

type FileType string

const (
	FileTypeGraph = FileType("sml")
	FileTypeLib   = FileType("alml")
	FileTypePlat  = FileType("spml")
	FileTypeMap   = FileType("mml")
)

var allFileTypes = [4]FileType{
	FileTypeGraph,
	FileTypeLib,
	FileTypePlat,
	FileTypeMap,
}

var descriptionFileTypes = map[FileType]string{
	FileTypeGraph: "Graph File (*.sml)",
	FileTypeLib:   "LibraryIf File (*.alml)",
	FileTypePlat:  "Platform File (*.spml)",
	FileTypeMap:   "Mapping File (*.mml)",
}

func Suffix(obj tr.ToplevelTreeElement) string {
	return string(fileType(obj))
}

func fileType(obj tr.ToplevelTreeElement) (ft FileType) {
	switch obj.(type) {
	case bh.SignalGraphIf:
		ft = FileTypeGraph
	case bh.LibraryIf:
		ft = FileTypeLib
	case pf.PlatformIf:
		ft = FileTypePlat
	case mp.MappingIf:
		ft = FileTypeMap
	default:
		log.Panicf("Suffix: invalid object type %T\n", obj)
	}
	return
}

func Description(obj tr.ToplevelTreeElement) string {
	return descriptionFileTypes[fileType(obj)]
}

type DirectoryMgr struct {
	xmlRoot string
	dirList []string
	current [4]string
}

var _ DirectoryMgrIf = (*DirectoryMgr)(nil)

func DirectoryMgrNew(xmlRoot string) (d *DirectoryMgr) {
	d = &DirectoryMgr{xmlRoot, nil, [4]string{xmlRoot, xmlRoot, xmlRoot, xmlRoot}}
	d.dirList = append(d.dirList, xmlRoot)
	return d
}

func (d DirectoryMgr) DirList() []string {
	return d.dirList
}

func (d DirectoryMgr) Current(ft FileType) (path string) {
	path = d.current[d.objTypeIndex(ft)]
	if !tool.IsSubPath("/", path) {
		path = fmt.Sprintf("%s/%s", d.xmlRoot, path)
	}
	return
}

func (d *DirectoryMgr) SetCurrent(ft FileType, newCurrent string) {
	for i, s := range d.dirList {
		if s == newCurrent {
			if i == 0 {
				return
			}
			d.current[d.objTypeIndex(ft)] = s
			tmpList := []string{s}
			for j := 0; j < i; j++ {
				tmpList = append(tmpList, d.dirList[j])
			}
			for j := i + 1; j < len(d.dirList); j++ {
				tmpList = append(tmpList, d.dirList[j])
			}
			d.dirList = tmpList
			return
		}
	}
	d.current[d.objTypeIndex(ft)] = newCurrent
	tmpList := []string{newCurrent}
	for j := 0; j < len(d.dirList); j++ {
		tmpList = append(tmpList, d.dirList[j])
	}
	d.dirList = tmpList
	log.Printf("DirectoryMgr.SetCurrent(%s, %s): %s\n", string(ft), newCurrent, d.Current(ft))
}

func (d *DirectoryMgr) FilenameToShow(filepath string) (prefix, filename string) {
	// TODO: evaluate FREESP_SEARCH_PATH
	if tool.IsSubPath(d.xmlRoot, filepath) {
		filename = tool.RelPath(d.xmlRoot, filepath)
		prefix = d.xmlRoot
	} else {
		filename = filepath
	}
	return
}

func (d *DirectoryMgr) objTypeIndex(ft FileType) int {
	switch ft {
	case FileTypeGraph:
		return 0
	case FileTypeLib:
		return 1
	case FileTypePlat:
		return 2
	case FileTypeMap:
		return 3
	default:
		log.Panicf("DirectoryMgr.objTypeIndex: invalid file type %s\n", string(ft))
	}
	return -1
}

var (
	dirMgr DirectoryMgrIf
)

func MenuFileInit(menu *GoAppMenu) {
	dirMgr = DirectoryMgrNew(backend.XmlRoot())
	menu.fileNewSg.Connect("activate", func() { fileNewSg(global.fts, global.ftv) })
	menu.fileNewLib.Connect("activate", func() { fileNewLib(global.fts, global.ftv) })
	menu.fileNewPlat.Connect("activate", func() { fileNewPlat(global.fts, global.ftv) })
	menu.fileNewMap.Connect("activate", func() { fileNewMap(global.fts, global.ftv) })
	menu.fileOpen.Connect("activate", func() { fileOpen(global.fts, global.ftv) })
	menu.fileSave.Connect("activate", func() { fileSave(global.fts) })
	menu.fileSaveAs.Connect("activate", func() { fileSaveAs(global.fts) })
	menu.fileClose.Connect("activate", func() { fileClose(menu, global.fts, global.ftv, global.jl) })

}

/*
 *		Callbacks
 */

func fileNewSg(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.SignalGraphMgr().New()
	if err != nil {
		log.Printf("fileNewSg: %s\n", err)
	}
}

func fileNewLib(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.LibraryMgr().New()
	if err != nil {
		log.Printf("fileNewLib: %s\n", err)
	}
}

func fileNewPlat(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.PlatformMgr().New()
	if err != nil {
		log.Printf("fileNewPlat: %s\n", err)
	}
}

func fileNewMap(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	var sgname, pfname string
	var ok bool
	sgname, ok = runFileDialog([]FileType{FileTypeGraph}, false, "Choose Signal Graph")
	if !ok {
		return
	}
	pfname, ok = runFileDialog([]FileType{FileTypePlat}, false, "Choose Platform")
	if !ok {
		return
	}
	var f tr.ToplevelTreeElement
	var err error
	prefix, fname := dirMgr.FilenameToShow(sgname)
	f, err = global.SignalGraphMgr().Access(fname)
	if err != nil {
		log.Printf("fileNewMap: %s\n", err)
		return
	}
	f.SetPathPrefix(prefix)
	global.MappingMgr().SetGraphForNew(f.(bh.SignalGraphIf))
	prefix, fname = dirMgr.FilenameToShow(pfname)
	f, err = global.PlatformMgr().Access(fname)
	if err != nil {
		log.Printf("fileNewMap: %s\n", err)
		return
	}
	f.SetPathPrefix(prefix)
	global.MappingMgr().SetPlatformForNew(f.(pf.PlatformIf))
	_, err = global.MappingMgr().New()
	if err != nil {
		log.Printf("fileNewMap: %s\n", err)
		return
	}
}

func getFileMgr(ft FileType) (fileMgr mod.FileManagerIf, err error) {
	switch ft {
	case FileTypeGraph:
		fileMgr = global.SignalGraphMgr()
	case FileTypeLib:
		fileMgr = global.LibraryMgr()
	case FileTypePlat:
		fileMgr = global.PlatformMgr()
	case FileTypeMap:
		fileMgr = global.MappingMgr()
	default:
		err = fmt.Errorf("fileMgr error: invalid file type %s\n", string(ft))
	}
	return
}

func fileOpen(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	log.Println("fileOpen")
	filename, ok := getFilenameToOpen()
	if !ok {
		return
	}
	prefix, fname := dirMgr.FilenameToShow(filename)
	var err error
	var obj tr.ToplevelTreeElement
	var fileMgr mod.FileManagerIf
	fileMgr, err = getFileMgr(FileType(tool.Suffix(fname)))
	obj, err = fileMgr.Access(fname)
	if err != nil {
		log.Printf("fileOpen error: %s\n", err)
		return
	}
	obj.SetPathPrefix(prefix)
	dirMgr.SetCurrent(fileType(obj), tool.Dirname(fname))
}

func fileSaveAs(fts *models.FilesTreeStore) {
	log.Println("fileSaveAs")
	filename, ok := getFilenameToSave(fts)
	if !ok {
		return
	}
	obj := getCurrentTopObject(fts)
	oldName := obj.Filename()
	prefix, fname := dirMgr.FilenameToShow(filename)
	err := global.FileMgr(obj).Rename(oldName, fname)
	if err != nil {
		log.Printf("fileSaveAs: %s\n", err)
		return
	}
	obj.SetPathPrefix(prefix)
	err = global.FileMgr(obj).Store(fname)
	if err != nil {
		log.Printf("fileSaveAs: %s\n", err)
		return
	}
	dirMgr.SetCurrent(fileType(obj), tool.Dirname(fname))
}

func isGeneratedFilename(filename string) bool {
	return strings.HasPrefix(filename, "new-file-") && !strings.Contains(filename, "/")
}

func fileSave(fts *models.FilesTreeStore) {
	log.Println("fileSave")
	obj := getCurrentTopObject(fts)
	filename := obj.Filename()
	if isGeneratedFilename(filename) {
		oldName := filename
		ok := true
		filename, ok = getFilenameToSave(fts)
		if !ok {
			return
		}
		dirMgr.SetCurrent(fileType(obj), tool.Dirname(filename))
		prefix, fname := dirMgr.FilenameToShow(filename)
		err := global.FileMgr(obj).Rename(oldName, fname)
		if err != nil {
			log.Printf("fileSave: %s\n", err)
			return
		}
		obj.SetPathPrefix(prefix)
	}
	dirMgr.SetCurrent(fileType(obj), tool.Dirname(filename))
	err := global.FileMgr(obj).Store(obj.Filename())
	if err != nil {
		log.Println(err)
	}
	return
}

func fileClose(menu *GoAppMenu, fts *models.FilesTreeStore, ftv *views.FilesTreeView, jl IJobList) {
	path := fts.GetCurrentId()
	if strings.Contains(path, ":") {
		path = strings.Split(path, ":")[0]
	}
	obj, err := fts.GetObjectById(path)
	if err != nil {
		return
	}
	global.FileMgr(obj).Remove(obj.(fd.Filenamer).Filename())
	jl.Reset()
	MenuEditPost(menu, fts, jl)
}

/*
 *		Local functions
 */

// empty fTypes means all types enabled
func runFileDialog(fTypes []FileType, toSave bool, announce string) (filename string, ok bool) {
	var action gtk.FileChooserAction
	var buttonText string
	if toSave {
		action = gtk.FILE_CHOOSER_ACTION_SAVE
		buttonText = "Save"
	} else {
		action = gtk.FILE_CHOOSER_ACTION_OPEN
		buttonText = "Open"
	}
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		announce,
		nil,
		action,
		"Cancel",
		gtk.RESPONSE_CANCEL,
		buttonText,
		gtk.RESPONSE_OK)
	if err != nil {
		log.Printf("runFileDialog error: %s\n", err)
		return
	}
	if len(fTypes) == 0 {
		for _, t := range allFileTypes {
			fTypes = append(fTypes, t)
		}
	}
	if len(fTypes) == 1 {
		log.Printf("runFileDialog: set current to %s.\n", dirMgr.Current(fTypes[0]))
		dialog.SetCurrentFolder(dirMgr.Current(fTypes[0]))
	} else {
		log.Printf("runFileDialog: set current to %s.\n", backend.XmlRoot())
		dialog.SetCurrentFolder(backend.XmlRoot())
	}
	for _, ft := range fTypes {
		ff, _ := gtk.FileFilterNew()
		ff.SetName(descriptionFileTypes[ft])
		ff.AddPattern(fmt.Sprintf("*.%s", string(ft)))
		dialog.AddFilter(ff)
	}
	response := dialog.Run()
	ok = (gtk.ResponseType(response) == gtk.RESPONSE_OK)
	filename = dialog.GetFilename()
	dialog.Destroy()
	return
}

func getFilenameToSave(fts *models.FilesTreeStore) (filename string, ok bool) {
	obj := getCurrentTopObject(fts)
	filename, ok = runFileDialog([]FileType{fileType(obj)}, true, "Save file as")
	// force correct suffix:
	dirname := tool.Dirname(filename)
	basename := tool.Basename(filename)
	suffix := tool.Suffix(basename)
	if suffix == basename { // no suffix given
		filename = fmt.Sprintf("%s.%s", filename, string(fileType(obj)))
	} else if suffix != string(fileType(obj)) { // wrong suffix
		log.Printf("getFilenameToSave: wrong suffix %s, forcing correct %s\n", suffix, string(fileType(obj)))
		idx := strings.LastIndex(basename, ".")
		if idx < 0 { // should never happen
			filename = fmt.Sprintf("%s/%s.%s", dirname, basename, string(fileType(obj)))
		} else {
			filename = fmt.Sprintf("%s/%s.%s", dirname, basename[:idx], string(fileType(obj)))
		}
	}
	return
}

func getFilenameToOpen() (filename string, ok bool) {
	return runFileDialog(nil, false, "Choose file to open")
}

func getToplevelId(fts *models.FilesTreeStore) string {
	p := fts.GetCurrentId()
	if p == "" {
		log.Fatal("fileSave error: fts.GetCurrentId() failed")
	}
	spl := strings.Split(p, ":")
	return spl[0] // TODO: move to function in fts...
}

func setCurrentTopValue(fts *models.FilesTreeStore, value string) {
	id0 := getToplevelId(fts)
	fts.SetValueById(id0, value)
}

func getCurrentTopObject(fts *models.FilesTreeStore) tr.ToplevelTreeElement {
	id0 := getToplevelId(fts)
	obj, err := fts.GetObjectById(id0)
	if err != nil {
		log.Fatal("fileSave error: fts.GetObjectByPath failed:", err)
	}
	return obj.(tr.ToplevelTreeElement)
}
