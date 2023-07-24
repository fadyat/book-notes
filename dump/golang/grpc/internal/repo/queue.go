package repo

type Queue[T any] interface {

	// Push pushes an element to the queue.
	Push(element *T)

	// Pop pops an element from the queue.
	Pop() *T

	// Peek returns the first element in the queue without removing it.
	Peek() *T

	// IsEmpty returns true if the queue is empty.
	IsEmpty() bool

	// Len returns the number of elements in the queue.
	Len() int
}

type node struct {
	value *Message
	next  *node
}

type queue struct {
	elements []*node
}

func (q *queue) Push(element *Message) {
	newNode := &node{value: element}
	if !q.IsEmpty() {
		q.elements[len(q.elements)-1].next = newNode
	}

	q.elements = append(q.elements, newNode)
}

func (q *queue) Pop() *Message {
	if q.IsEmpty() {
		return nil
	}

	first := q.elements[0]
	q.elements = q.elements[1:]
	return first.value
}

func (q *queue) Peek() *Message {
	if q.IsEmpty() {
		return nil
	}

	return q.elements[0].value
}

func (q *queue) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *queue) Len() int {
	return len(q.elements)
}
