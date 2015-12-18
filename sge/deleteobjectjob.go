package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
)

type DeleteObjectJob struct {
	deletedObjects []freesp.IdWithObject
	id             string
}

func DeleteObjectJobNew(id string) *DeleteObjectJob {
	return &DeleteObjectJob{nil, id}
}

func (j *DeleteObjectJob) String() string {
	ret := fmt.Sprintf("Delete %s\n", j.id)
	for _, obj := range j.deletedObjects {
		ret = fmt.Sprintf("  %v", obj)
	}
	return ret
}
