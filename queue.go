package concurrentcache

import (
	"errors"
	"sync"
)

// Queue struct
type Queue struct {
	head *QNode
	tail *QNode
	pool *sync.Pool
}

// QNode struct
type QNode struct {
	prev *QNode
	next *QNode
	V    interface{}
}

// NewQueue return a Queue
func NewQueue() *Queue {
	pool := &sync.Pool{
		New: func() interface{} {
			return &QNode{}
		},
	}
	return &Queue{pool: pool}
}

// NewQNode return a QNode
func (q *Queue) NewQNode(v interface{}) *QNode {
	n := q.pool.Get().(*QNode)
	n.reset()
	n.V = v
	return n
}

// Recycle reuse QNode
func (q *Queue) Recycle(n *QNode) {
	if n != nil {
		q.pool.Put(n)
	}
}

func (n *QNode) reset() {
	if n != nil {
		n.prev = nil
		n.next = nil
		n.V = nil
	}
}

// LPush a QNode to the left of the Queue
func (q *Queue) LPush(n *QNode) {
	q.push(n, true)
}

// RPush a QNode to the right of the Queue
func (q *Queue) RPush(n *QNode) {
	q.push(n, false)
}

func (q *Queue) push(n *QNode, left bool) {
	if q.Empty() {
		n.next, n.prev = nil, nil
		q.head, q.tail = n, n
		return
	}
	if left {
		n.next, n.prev = q.head, nil
		q.head.prev = n
		q.head = n
	} else {
		n.next, n.prev = nil, q.tail
		q.tail.next = n
		q.tail = n
	}
}

// InsertAfter prev QNode
func (q *Queue) InsertAfter(prev, n *QNode) (bool, error) {
	if prev == nil || n == nil {
		return false, errors.New("invalid node")
	}
	n.prev = prev
	n.next = prev.next
	if prev.next != nil {
		prev.next.prev = n
	}
	prev.next = n
	if prev == q.tail {
		q.tail = n
	}
	return true, nil
}

// InsertBefore next QNode
func (q *Queue) InsertBefore(next, n *QNode) (bool, error) {
	if next == nil || n == nil {
		return false, errors.New("invalid node")
	}
	n.next = next
	n.prev = next.prev
	if next.prev != nil {
		next.prev.next = n
	}
	next.prev = n
	if next == q.head {
		q.head = n
	}
	return true, nil
}

// LPop QNode from the left of Queue
func (q *Queue) LPop() *QNode {
	return q.pop(true)
}

// RPop QNode from the right of Queue
func (q *Queue) RPop() *QNode {
	return q.pop(false)
}

func (q *Queue) pop(left bool) *QNode {
	if q.Empty() {
		return nil
	}
	if left {
		n := q.head
		if q.head == q.tail {
			q.head, q.tail = nil, nil
		} else {
			q.head = q.head.next
		}
		n.next, n.prev = nil, nil
		return n
	}
	n := q.tail
	if q.head == q.tail {
		q.head, q.tail = nil, nil
	} else {
		q.tail = q.tail.prev
	}
	n.next, n.prev = nil, nil
	return n
}

// Delete a QNode
func (q *Queue) Delete(n *QNode) (bool, error) {
	if n == nil {
		return false, errors.New("invalid node")
	}
	if n.prev != nil {
		n.prev.next = n.next
	} else if n == q.head && n != nil {
		q.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else if n == q.tail && n != nil {
		q.tail = n.prev
	}
	n.prev = nil
	n.next = nil
	return true, nil
}

// Empty if Queue is empty
func (q *Queue) Empty() bool {
	if q.head == q.tail && q.head == nil {
		return true
	}
	return false
}
