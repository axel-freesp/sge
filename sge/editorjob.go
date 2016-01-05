package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"log"
)

type JobType int

const (
	JobNewElement JobType = iota
	JobDeleteObject
	JobEdit
)

type EditorJob struct {
	jobType      JobType
	jobDetail    fmt.Stringer
	newElement   *NewElementJob
	deleteObject *DeleteObjectJob
	edit         *EditJob
}

func EditorJobNew(jobType JobType, jobDetail fmt.Stringer) *EditorJob {
	ret := &EditorJob{jobType, jobDetail, nil, nil, nil}
	switch jobType {
	case JobNewElement:
		ret.newElement = jobDetail.(*NewElementJob)
	case JobDeleteObject:
		ret.deleteObject = jobDetail.(*DeleteObjectJob)
	case JobEdit:
		ret.edit = jobDetail.(*EditJob)
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

func (a *jobApplier) Apply(jobI interface{}) (state interface{}, err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		var object freesp.TreeElement
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
	}
	return
}

func (a *jobApplier) Revert(jobI interface{}) (state interface{}, err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		var del []freesp.IdWithObject
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
	}
	return
}
