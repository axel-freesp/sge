package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"log"
)

type fileManagerSG struct {
	filenameFactory
	context        FilemanagerContextIf
	signalGraphMap map[string]bh.SignalGraphIf
}

var _ mod.FileManagerIf = (*fileManagerSG)(nil)

func FileManagerSGNew(context FilemanagerContextIf) *fileManagerSG {
	return &fileManagerSG{FilenameFactoryInit("sml"), context, make(map[string]bh.SignalGraphIf)}
}

//
//      FileManagerIf interface
//

func (f *fileManagerSG) New() (sg tr.ToplevelTreeElement, err error) {
	filename := f.NewFilename()
	sg = behaviour.SignalGraphNew(filename, f.context)
	f.signalGraphMap[filename] = sg.(bh.SignalGraphIf)
	var newId string
	newId, err = f.context.FTS().AddToplevel(sg.(bh.SignalGraphIf))
	if err != nil {
		err = fmt.Errorf("fileManagerSG.New: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)
	var gv views.GraphViewIf
	gv, err = views.SignalGraphViewNew(sg.(bh.SignalGraphIf), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerSG.New: %s", err)
		return
	}
	f.context.GVC().Add(gv, filename)
	f.context.ShowAll()
	return
}

func (f *fileManagerSG) Access(name string) (sg tr.ToplevelTreeElement, err error) {
	var ok bool
	sg, ok = f.signalGraphMap[name]
	if ok {
		return
	}
	sg = behaviour.SignalGraphNew(name, f.context)
	var filedir string
	for _, filedir = range backend.XmlSearchPaths() {
		err = sg.ReadFile(fmt.Sprintf("%s/%s", filedir, name))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("fileManagerSG.Access: graph file %s not found.", name)
		return
	}
	var filename string
	if len(filedir) > 0 {
		filename = fmt.Sprintf("%s/%s", filedir, name)
	} else {
		filename = name
	}
	sg.SetPathPrefix(filedir)
	hint := backend.XmlGraphHintNew(name)
	hintfilename := f.HintFilename(filename)
	var buf []byte
	buf, err = tool.ReadFile(hintfilename)
	if err == nil {
		_, err = hint.Read(buf)
		if err != nil {
			return
		}
		err = behaviour.SignalGraphApplyHints(sg.(bh.SignalGraphIf), hint)
		if err != nil {
			err = fmt.Errorf("fileManagerSG.Access: %s", err)
			return
		}
	}
	var newId string
	newId, err = f.context.FTS().AddToplevel(sg.(bh.SignalGraphIf))
	if err != nil {
		err = fmt.Errorf("fileManagerSG.Access: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)

	var gv views.GraphViewIf
	gv, err = views.SignalGraphViewNew(sg.(bh.SignalGraphIf), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerSG.Access: %s", err)
		return
	}
	f.context.GVC().Add(gv, name)
	f.signalGraphMap[name] = sg.(bh.SignalGraphIf)
	log.Printf("fileManagerSG.Access: graph %s successfully loaded.\n", name)

	return
}

func (f *fileManagerSG) Remove(name string) {
	// TODO: check if used in existing mappings
	sg, ok := f.signalGraphMap[name]
	if !ok {
		log.Printf("fileManagerSG.Remove warning: graph %s not found.\n", name)
		return
	}
	delete(f.signalGraphMap, name)
	var nodes []bh.NodeIf
	for _, n := range sg.ItsType().Nodes() {
		nodes = append(nodes, n)
	}
	f.context.CleanupNodeTypesFromNodes(nodes)
	f.context.CleanupSignalTypesFromNodes(nodes)
	id, _ := f.context.FTS().GetToplevelId(sg)
	f.context.FTS().RemoveToplevel(id)
	f.context.GVC().RemoveGraphView(sg)
}

func (f *fileManagerSG) Rename(oldName, newName string) (err error) {
	sg, ok := f.signalGraphMap[oldName]
	if !ok {
		err = fmt.Errorf("fileManagerSG.Rename error: graph %s not found\n", oldName)
		return
	}
	_, ok = f.signalGraphMap[newName]
	if ok {
		err = fmt.Errorf("fileManagerSG.Rename error: cannot rename graph %s to be %s: already exists\n", oldName, newName)
		return
	}
	id, _ := f.context.FTS().GetToplevelId(sg)
	f.context.FTS().SetValueById(id, newName)
	f.context.GVC().Rename(oldName, newName)
	delete(f.signalGraphMap, oldName)
	sg.SetFilename(newName)
	f.signalGraphMap[newName] = sg
	return
}

func (f *fileManagerSG) Store(name string) (err error) {
	sg, ok := f.signalGraphMap[name]
	if !ok {
		err = fmt.Errorf("fileManagerSG.Rename error: graph %s not found\n", name)
		return
	}
	var filename string
	if len(sg.PathPrefix()) == 0 {
		filename = sg.Filename()
	} else {
		filename = fmt.Sprintf("%s/%s", sg.PathPrefix(), sg.Filename())
	}
	err = sg.WriteFile(filename)
	if err != nil {
		return
	}
	hint := behaviour.CreateXmlGraphHint(sg)
	hintfilename := f.HintFilename(filename)
	var buf []byte
	buf, err = hint.Write()
	if err != nil {
		return
	}
	err = tool.WriteFile(hintfilename, buf)
	return
}
