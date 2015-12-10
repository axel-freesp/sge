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
)

type EditorJob struct {
	jobType    JobType
	jobDetail  fmt.Stringer
	newElement *NewElementJob
}

func EditorJobNew(jobType JobType, jobDetail fmt.Stringer) *EditorJob {
	ret := &EditorJob{jobType, jobDetail, nil}
	switch jobType {
	case JobNewElement:
		ret.newElement = jobDetail.(*NewElementJob)
	}
	return ret
}

func (e *EditorJob) String() string {
	var kind string
	switch e.jobType {
	case JobNewElement:
		kind = "NewElement"
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
		switch object.(type) {
		case freesp.Implementation, freesp.SignalType, freesp.NamedPortType, freesp.NodeType, freesp.Node, freesp.Port:
			job.newElement.newId, err = a.fts.AddNewObject(job.newElement.parentId, object)
		default:
			log.Println("jobApplier.Apply(JobNewElement): invalid data type")
		}
	}
	return
}

func (j *jobApplier) Revert(jobI interface{}) error {
	return nil
}
