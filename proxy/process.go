package proxy

import "time"

//Process the process flow
type Process struct {
	id           string //process call id
	caller       Conn   //caller conn
	callee       Conn   //callee conn
	route        *Route //process route
	relationship *Relationship
	stacks       []*Message //stacks
	createdAt    time.Time
	updatedAt    time.Time
}

func (proc *Process) Ready() bool {
	return proc.caller != nil && proc.callee != nil
}

func (proc *Process) Caller() Conn {
	return proc.caller
}

func (proc *Process) Callee() Conn {
	return proc.callee
}

func (proc *Process) Push(msg *Message) {
	proc.updatedAt = time.Now()
	proc.stacks = append(proc.stacks, msg)
}

func NewProcess(id string) *Process {
	return &Process{id: id, createdAt: time.Now(), stacks: make([]*Message, 0)}
}
