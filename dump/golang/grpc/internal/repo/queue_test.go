package repo

import (
	"bytes"
	"testing"
)

func compareMessages(t *testing.T, expected, actual *Message) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil || actual == nil {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	if expected.offset != actual.offset {
		t.Errorf("expected %v, got %v", expected.offset, actual.offset)
	}

	if !bytes.Equal(expected.content, actual.content) {
		t.Errorf("expected %v, got %v", string(expected.content), string(actual.content))
	}
}

func TestQueue_PushAndPop(t *testing.T) {
	testCases := []struct {
		name     string
		messages []*Message
		popCount int
	}{
		{
			name: "success, pop none",
			messages: []*Message{
				{offset: 0, content: []byte("a")},
				{offset: 1, content: []byte("b")},
				{offset: 2, content: []byte("c")},
			},
		},
		{
			name: "success, pop one",
			messages: []*Message{
				{offset: 0, content: []byte("a")},
				{offset: 1, content: []byte("b")},
				{offset: 2, content: []byte("c")},
			},
			popCount: 1,
		},
		{
			name:     "success, pop from empty queue",
			popCount: 1,
		},
		{
			name: "success, pop all",
			messages: []*Message{
				{offset: 0, content: []byte("a")},
				{offset: 1, content: []byte("b")},
			},
			popCount: 2,
		},
	}

	checkLength := func(t *testing.T, expected, actual int) {
		if expected != actual {
			t.Errorf("expected %d, got %d", expected, actual)
		}
	}

	checkForNextMessage := func(t *testing.T, idx int, q *queue) {
		if idx < q.Len()-1 {
			compareMessages(t, q.elements[idx+1].value, q.elements[idx].next.value)
			return
		}

		if q.elements[idx].next != nil {
			t.Errorf("expected %v, got %v", nil, q.elements[idx].next)
		}
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}

		return b
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := &queue{}
			for _, m := range tc.messages {
				q.Push(m)
			}

			checkLength(t, len(tc.messages), q.Len())
			for i, m := range tc.messages {
				compareMessages(t, m, q.elements[i].value)
				checkForNextMessage(t, i, q)
			}

			for i := 0; i < tc.popCount; i++ {
				q.Pop()
			}

			checkLength(t, max(len(tc.messages)-tc.popCount, 0), q.Len())
			for i := 0; i < q.Len(); i++ {
				compareMessages(t, tc.messages[i+tc.popCount], q.elements[i].value)
				checkForNextMessage(t, i, q)
			}
		})
	}
}

func TestQueue_Peek(t *testing.T) {
	q := &queue{}
	top := q.Peek()
	compareMessages(t, nil, top)

	msg := &Message{offset: 0, content: []byte("a")}
	q.Push(msg)
	top = q.Peek()
	compareMessages(t, msg, top)
}
