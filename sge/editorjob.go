package main

import (
	"fmt"
	//"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"log"
)

type JobType int

const (
	JobNewElement JobType = iota
	JobDeleteObject
)

type EditorJob struct {
	jobType      JobType
	jobDetail    fmt.Stringer
	newElement   *NewElementJob
	deleteObject *DeleteObjectJob
}

func EditorJobNew(jobType JobType, jobDetail fmt.Stringer) *EditorJob {
	ret := &EditorJob{jobType, jobDetail, nil, nil}
	switch jobType {
	case JobNewElement:
		ret.newElement = jobDetail.(*NewElementJob)
	case JobDeleteObject:
		ret.deleteObject = jobDetail.(*DeleteObjectJob)
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
	}
	return fmt.Sprintf("EditorJob( %s( %v ) )", kind, e.jobDetail)
}

type jobApplier struct {
	// IJobApplier
	fts *models.FilesTreeStore
}

func jobApplierNew(fts *models.FilesTreeStore) *jobApplier {
	j := &jobApplier{fts}

	return j
}

func (a *jobApplier) Apply(jobI interface{}) (err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		object := job.newElement.CreateObject(a.fts)
		job.newElement.newId, err = a.fts.AddNewObject(job.newElement.parentId, -1, object)
	case JobDeleteObject:
		job.deleteObject.deletedObjects, err = a.fts.DeleteObject(job.deleteObject.id)
	}
	return
}

func (a *jobApplier) Revert(jobI interface{}) (err error) {
	job := jobI.(*EditorJob)
	switch job.jobType {
	case JobNewElement:
		log.Println("jobApplier.Revert: delete added object at ", job.newElement.newId)
		_, err = a.fts.DeleteObject(job.newElement.newId)
		if err != nil {
			log.Println(err)
		}
	case JobDeleteObject:
		var i int
		for i = len(job.deleteObject.deletedObjects) - 1; i >= 0; i-- {
			d := job.deleteObject.deletedObjects[i]
			log.Printf("jobApplier.Revert: adding deleted object %v at %s/%d\n", d.Object, d.ParentId, d.Position)
			_, err = a.fts.AddNewObject(d.ParentId, d.Position, d.Object)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return
}
