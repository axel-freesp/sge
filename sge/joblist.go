package main

import (
	"github.com/axel-freesp/sge/tool"
	"log"
)

const jobListCapacity = 10

type IJobApplier interface {
	Apply(job interface{}) (state interface{}, err error)
	Revert(job interface{}) (state interface{}, err error)
}

type IJobList interface {
	Undo() (state interface{}, ok bool)
	Redo() (state interface{}, ok bool)
	CanUndo() bool
	CanRedo() bool
	Apply(job interface{}) (state interface{}, ok bool)
}

type jobList struct {
	undoStack, redoStack *tool.Stack
	system               IJobApplier
}

var _ IJobList = (*jobList)(nil)

func jobListNew(system IJobApplier) *jobList {
	j := &jobList{tool.StackNew(), tool.StackNew(), system}
	return j
}

func (j *jobList) reset() {
	j.undoStack.Reset()
	j.redoStack.Reset()
}

func (j *jobList) Undo() (state interface{}, ok bool) {
	ok = false
	if j.CanUndo() {
		job := j.undoStack.Pop()
		j.redoStack.Push(job)
		var err error
		state, err = j.system.Revert(job)
		if err != nil {
			log.Println("jobList.Undo error: ", err)
			j.reset()
		}
		ok = true
	}
	return
}

func (j *jobList) Redo() (state interface{}, ok bool) {
	ok = false
	if j.CanRedo() {
		job := j.redoStack.Pop()
		j.undoStack.Push(job)
		var err error
		state, err = j.system.Apply(job)
		if err != nil {
			log.Println("jobList.Redo error: ", err)
			j.reset()
		}
		ok = true
	}
	return
}

func (j *jobList) Apply(job interface{}) (state interface{}, ok bool) {
	ok = false
	state, err := j.system.Apply(job)
	j.redoStack.Reset()
	if err == nil {
		j.undoStack.Push(job)
		ok = true
	} else {
		log.Println("jobList.Apply error: ", err)
	}
	return
}

func (j *jobList) CanUndo() bool {
	return !j.undoStack.IsEmpty()
}

func (j *jobList) CanRedo() bool {
	return !j.redoStack.IsEmpty()
}
