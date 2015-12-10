package main

import (
	"log"
)

const jobListCapacity = 10

type IJobApplier interface {
	Apply(job interface{}) error
	Revert(job interface{}) error
}

type IJobList interface {
	Undo()
	Redo()
	CanUndo() bool
	CanRedo() bool
	Apply(job interface{}) (ok bool)
}

type jobList struct {
	// IJobList
	undoStack, redoStack []interface{}
	system               IJobApplier
}

func jobListNew(system IJobApplier) *jobList {
	j := &jobList{make([]interface{}, 0, jobListCapacity),
		make([]interface{}, 0, jobListCapacity),
		system}
	return j
}

func (j *jobList) Undo() {
	log.Println("jobList.Undo")
	undoLen := len(j.undoStack)
	if undoLen > 0 {
		job := j.undoStack[undoLen-1]
		j.redoStack = append(j.redoStack, job)
		j.undoStack = j.undoStack[:undoLen-1]
		j.system.Revert(job)
	}
}

func (j *jobList) Redo() {
	log.Println("jobList.Redo")
	redoLen := len(j.redoStack)
	if redoLen > 0 {
		job := j.redoStack[redoLen-1]
		err := j.system.Apply(job)
		if err == nil {
			j.undoStack = append(j.undoStack, job)
			j.redoStack = j.redoStack[:redoLen-1]
		} else {
			j.undoStack = j.undoStack[:0]
			j.redoStack = j.redoStack[:0]
		}
	}
}

func (j *jobList) Apply(job interface{}) (ok bool) {
	log.Println("jobList.Apply", job)
	err := j.system.Apply(job)
	j.redoStack = j.redoStack[:0]
	if err == nil {
		j.undoStack = append(j.undoStack, job)
	}
	return err == nil
}

func (j *jobList) CanUndo() bool {
	return true
}

func (j *jobList) CanRedo() bool {
	return len(j.redoStack) > 0
}
