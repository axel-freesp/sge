package main

import (
	"fmt"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	"log"
)

type JobType int

const (
	JobNewElement JobType = iota
	JobDeleteObject
	JobEdit
	JobPaste
)

type EditorJob struct {
	jobType      JobType
	jobDetail    fmt.Stringer
	newElement   *NewElementJob
	deleteObject *DeleteObjectJob
	edit         *EditJob
	paste        *PasteJob
}

func EditorJobNew(jobType JobType, jobDetail fmt.Stringer) *EditorJob {
	ret := &EditorJob{jobType, jobDetail, nil, nil, nil, nil}
	switch jobType {
	case JobNewElement:
		ret.newElement = jobDetail.(*NewElementJob)
	case JobDeleteObject:
		ret.deleteObject = jobDetail.(*DeleteObjectJob)
	case JobEdit:
		ret.edit = jobDetail.(*EditJob)
	case JobPaste:
		ret.paste = jobDetail.(*PasteJob)
	}
	return ret
}

func (e *EditorJob) String() string {
	var kind string
	switch e.jobType {
	case JobNewElement:
		kind = "NewElement"
	case JobDeleteObject:
		kind = "DeleteObject"
	case JobEdit:
		kind = "Edit"
	case JobPaste:
		kind = "Paste"
	}
	return fmt.Sprintf("EditorJob( %s( %v ) )", kind, e.jobDetail)
}

type jobApplier struct {
	fts *models.FilesTreeStore
}

var _ IJobApplier = (*jobApplier)(nil)

func jobApplierNew(fts *models.FilesTreeStore) *jobApplier {
	j := &jobApplier{fts}

	return j
}

var level int = 0

func (a *jobApplier) Apply(jobI interface{}) (state interface{}, err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		var object tr.TreeElementIf
		object, err = job.newElement.CreateObject(a.fts)
		if err != nil {
			log.Printf("jobApplier.Apply error (JobNewElement): %s\n", err)
			return
		}
		job.newElement.newId, err = a.fts.AddNewObject(job.newElement.parentId, -1, object)
		if err == nil {
			state = job.newElement.newId
		} else {
			state = a.fts.GetCurrentId()
			log.Printf("jobApplier.Apply error (JobNewElement): %s\n", err)
		}
	case JobDeleteObject:
		job.deleteObject.deletedObjects, err = a.fts.DeleteObject(job.deleteObject.id)
		if err == nil {
			d := job.deleteObject.deletedObjects[len(job.deleteObject.deletedObjects)-1]
			state = d.ParentId
		} else {
			state = a.fts.GetCurrentId()
			log.Printf("jobApplier.Apply (JobDeleteObject): error: %s\n", err)
		}
	case JobEdit:
		state, err = job.edit.EditObject(a.fts, EditJobForward)
		if err != nil {
			log.Printf("jobApplier.Apply (JobEdit): error: %s\n", err)
		}
	case JobPaste:
		level++
		for _, j := range job.paste.newElements {
			//log.Printf("jobApplier.Apply (JobPaste): level %d: %v\n", level, j)
			j.parentId = job.paste.context
			state, err = a.Apply(EditorJobNew(JobNewElement, j))
			if err != nil {
				log.Printf("jobApplier.Apply (JobPaste): error: %s\n", err)
				return
			}
			for _, ch := range job.paste.children {
				//log.Printf("jobApplier.Apply (JobPaste): level %d: %v\n", level, ch)
				ch.context = state.(string)
				_, err = a.Apply(EditorJobNew(JobPaste, ch))
				if err != nil {
					log.Printf("jobApplier.Apply (JobPaste): error: %s\n", err)
				}
			}
		}
		level--
	}
	return
}

func (a *jobApplier) Revert(jobI interface{}) (state interface{}, err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		var del []tr.IdWithObject
		del, err = a.fts.DeleteObject(job.newElement.newId)
		if err == nil {
			state = del[0].ParentId
		} else {
			state = a.fts.GetCurrentId()
			log.Printf("jobApplier.Revert (JobNewElement): error: %s\n", err)
		}
	case JobDeleteObject:
		d := job.deleteObject.deletedObjects[len(job.deleteObject.deletedObjects)-1]
		state = fmt.Sprintf("%s:%d", d.ParentId, d.Position)
		for i := len(job.deleteObject.deletedObjects) - 1; i >= 0; i-- {
			d = job.deleteObject.deletedObjects[i]
			_, err = a.fts.AddNewObject(d.ParentId, d.Position, d.Object)
			if err != nil {
				state = a.fts.GetCurrentId()
				log.Printf("jobApplier.Revert (JobDeleteObject): error: %s\n", err)
				return
			}
		}
	case JobEdit:
		state, err = job.edit.EditObject(a.fts, EditJobRevert)
		if err != nil {
			log.Printf("jobApplier.Revert (JobEdit): error: %s\n", err)
		}
	case JobPaste:
		state, err = a.Revert(EditorJobNew(JobNewElement, job.paste.newElements[0]))
		if err != nil {
			log.Printf("jobApplier.Apply (JobPaste): error: %s\n", err)
			return
		}
	}
	return
}
