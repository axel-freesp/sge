package tool

// TODO: synchronization

type Stack struct {
	stack []interface{}
}

func StackNew() *Stack {
	return &Stack{nil}
}

func (s Stack) IsEmpty() bool {
	return len(s.stack) == 0
}

func (s *Stack) Reset() {
	s.stack = s.stack[:0]
}

func (s *Stack) Push(n interface{}) {
	s.stack = append(s.stack, n)
}

func (s *Stack) Pop() (n interface{}) {
	index := len(s.stack) - 1
	n = s.stack[index]
	s.stack = s.stack[:index]
	return
}
