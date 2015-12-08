package main

import (
	"fmt"
)

type JobType int

const (
	JobNewElement JobType = iota
)

type EditorJob struct {
	jobType   JobType
	jobDetail fmt.Stringer
}

func EditorJobNew(jobType JobType, jobDetail fmt.Stringer) *EditorJob {
	return &EditorJob{jobType, jobDetail}
}

func (e *EditorJob) String() string {
	var kind string
	switch e.jobType {
	case JobNewElement:
		kind = "NewElement"
	}
	return fmt.Sprintf("EditorJob( %s( %v ) )", kind, e.jobDetail)
}
