package stack

type (
	Stack struct {
		top  *node
		size int
	}

	node struct {
		value interface{}
		prev  *node
	}
)

// New create new stack
func New() *Stack {
	return &Stack{nil, 0}
}

// Push new data to the stack
func (s *Stack) Push(value interface{}) {
	if value == nil || value == "" {
		return
	}

	n := &node{
		value: value,
		prev:  s.top,
	}

	s.top = n
	s.size++
}

// Pop data from the stack, return top value
func (s *Stack) Pop() (value interface{}) {
	if s.size == 0 || s.top == nil {
		return nil
	}

	n := s.top

	s.top = n.prev
	s.size--

	n.prev = nil

	return n.value
}

// Peek return top value from the stack without removing the top value
func (s Stack) Peek() (value interface{}) {
	if s.top == nil {
		return nil
	}

	return s.top.value
}

func (s Stack) Size() int {
	return s.size
}
