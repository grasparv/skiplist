package skiplist

type searchResult[T any] struct {
	// prev points to the node to the left of the result if it was
	// was not found, or at the top of the column
	prev *node[T]

	// result points to a node if a result was found
	result *node[T]

	// descends contains a history of each node that we walked down from on previous levels
	descends []*node[T]
}

func (r *searchResult[T]) wasFound() bool {
	return r.result != nil
}

type compareFunc[T any] func(a, b *T) bool

func find[T any](n *node[T], value *T, equalTo compareFunc[T], greaterThan compareFunc[T]) (r searchResult[T]) {
	r.descends = make([]*node[T], 0)

	for n != nil {
		// find the lowest node in the column
		bottomNode := n.bottomNode()

		if bottomNode.value != nil {
			// we are not in the head column
			if equalTo(bottomNode.value, value) {
				// the exact value was found
				r.prev = n
				r.result = bottomNode
				return
			} else if greaterThan(bottomNode.value, value) {
				// we passed the value in a lower level,
				// backtrack to the left and go down
				n = n.prev
				if n.down != nil {
					r.descends = append(r.descends, n)
					n = n.down
				} else {
					// there is no exact hit, we are to the left of its void
					r.prev = n
					return
				}
				continue
			}
		}

		r.prev = n
		if n.next != nil {
			// we have not passed the value yet,
			// advance to the right
			n = n.next
		} else if n.down != nil {
			// (only in head) we have no more right links, only down;
			// go as far down as possible to find a new right
			r.descends = append(r.descends, n)
			n = n.down
		} else {
			return
		}
	}

	return
}
