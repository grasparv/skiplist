package skiplist

import (
	"math/rand"
	"sync"
)

type Skiplist[T any] struct {
	levels      []*node[T]
	greaterThan compareFunc[T]
	equalTo     compareFunc[T]
	count       int
	lock        sync.RWMutex
}

type node[T any] struct {
	next  *node[T]
	prev  *node[T]
	down  *node[T]
	value *T
}

func (n *node[T]) bottomNode() *node[T] {
	for n.down != nil {
		n = n.down
	}
	return n
}

func New[T any](levels int, equalTo compareFunc[T], greaterThan compareFunc[T]) *Skiplist[T] {
	if levels < 1 {
		panic("too few levels")
	}

	l := &Skiplist[T]{
		levels:      make([]*node[T], levels),
		greaterThan: greaterThan,
		equalTo:     equalTo,
	}

	var prev *node[T]
	for i := 0; i < levels; i++ {
		l.levels[i] = &node[T]{
			down: prev,
		}
		prev = l.levels[i]
	}

	return l
}

func (l *Skiplist[T]) getTopNode() *node[T] {
	return l.levels[len(l.levels)-1]
}

func (l *Skiplist[T]) Insert(value *T) {
	l.lock.Lock()
	defer l.lock.Unlock()

	n := l.getTopNode()

	r := find(n, value, l.equalTo, l.greaterThan)
	if r.wasFound() {
		// only allow unique values in list
		return
	}
	l.count++

	// insert the node on the first level
	newValueNode := &node[T]{
		value: value,
	}
	prevRight := r.prev.next
	l.link(r.prev, newValueNode)
	if prevRight != nil {
		l.link(newValueNode, prevRight)
	}

	// build a stack on top of the new value node
	for n, i, walkdowns := newValueNode, 1, len(r.descends)-1; l.coinToss() && i < len(l.levels) && walkdowns >= 0; i++ {
		// 50% chance to populate the node we came down from's right neighbour
		descendNode := r.descends[walkdowns]
		descendRight := descendNode.next
		newNode := &node[T]{
			down: n,
		}
		// link node we came down from's right with the new node and connect back
		l.link(descendNode, newNode)
		if descendRight != nil {
			// link the former right node with the new node
			l.link(newNode, descendRight)
		}
		n = newNode
		walkdowns--
	}
}

func (l *Skiplist[T]) coinToss() bool {
	return rand.Intn(2) == 0
}

func (l *Skiplist[T]) link(a, b *node[T]) {
	a.next = b
	b.prev = a
}

func (l *Skiplist[T]) unlink(a *node[T]) {
	left := a.prev
	right := a.next
	if left != nil && right != nil {
		left.next = right
		right.prev = left
	} else if left != nil && right == nil {
		left.next = nil
	} else if left == nil && right != nil {
		right.prev = nil
	}
}

func (l *Skiplist[T]) Contains(value *T) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	n := l.getTopNode()
	r := find(n, value, l.equalTo, l.greaterThan)
	return r.wasFound()
}

func (l *Skiplist[T]) Remove(value *T) {
	l.lock.Lock()
	defer l.lock.Unlock()

	n := l.getTopNode()

	r := find(n, value, l.equalTo, l.greaterThan)
	if !r.wasFound() {
		// only allow unique values in list
		return
	}
	l.count--

	l.unlink(r.result)
	for n := r.prev; n != nil; n = n.down {
		l.unlink(n)
	}
}

func (l *Skiplist[T]) Length() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.count
}

func (l *Skiplist[T]) Range(startingPoint *T, callback func(i int, value *T) bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	n := l.getTopNode()

	r := find(n, startingPoint, l.equalTo, l.greaterThan)
	if r.wasFound() {
		n = r.result
	} else {
		n = r.prev
	}

	for i := 0; n != nil; n = n.next {
		x := n.bottomNode().value
		if x == nil || l.greaterThan(startingPoint, x) {
			continue
		}
		if !callback(i, x) {
			break
		}
	}
}
