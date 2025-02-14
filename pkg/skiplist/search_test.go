package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func pointerTo[T any](v T) *T {
	return &v
}

func eqInt(a, b *int) bool {
	return *a == *b
}

func gtInt(a, b *int) bool {
	return *a > *b
}

func buildFixture[T any]() *Skiplist[int] {
	l := New(7, eqInt, gtInt)

	///////////////////////////////////////////////////
	// LEVEL 0
	///////////////////////////////////////////////////

	level0four := &node[int]{
		value: pointerTo(4),
	}

	level0six := &node[int]{
		value: pointerTo(6),
	}

	level0nine := &node[int]{
		value: pointerTo(9),
	}

	level0seventeen := &node[int]{
		value: pointerTo(17),
	}

	level0twentytwo := &node[int]{
		value: pointerTo(22),
	}

	level0twentynine := &node[int]{
		value: pointerTo(29),
	}

	l.levels[0].next = level0four
	level0four.next = level0six
	level0six.next = level0nine
	level0nine.next = level0seventeen
	level0seventeen.next = level0twentytwo
	level0twentytwo.next = level0twentynine

	level0four.prev = l.levels[0]
	level0six.prev = level0four
	level0nine.prev = level0six
	level0seventeen.prev = level0nine
	level0twentytwo.prev = level0seventeen
	level0twentynine.prev = level0twentytwo

	///////////////////////////////////////////////////
	// LEVEL 1
	///////////////////////////////////////////////////

	level10 := &node[int]{
		down: level0four,
	}

	level11 := &node[int]{
		down: level0nine,
	}

	level12 := &node[int]{
		down: level0twentytwo,
	}

	l.levels[1].next = level10
	level10.next = level11
	level11.next = level12

	level10.prev = l.levels[1]
	level11.prev = level10
	level12.prev = level11

	///////////////////////////////////////////////////
	// LEVEL 2
	///////////////////////////////////////////////////

	level20 := &node[int]{
		down: level11,
	}

	level21 := &node[int]{
		down: level12,
	}

	l.levels[2].next = level20
	level20.next = level21
	level20.prev = l.levels[2]
	level21.prev = level20

	///////////////////////////////////////////////////
	// LEVEL 3
	///////////////////////////////////////////////////

	level30 := &node[int]{
		down: level21,
	}

	l.levels[3].next = level30
	level30.prev = l.levels[3]

	return l
}

func TestSearch(t *testing.T) {
	l := buildFixture[int]()

	t.Run("miss 15", func(t *testing.T) {
		n := l.getTopNode()
		r := find(n, pointerTo(15), eqInt, gtInt)
		require.False(t, r.wasFound())
		require.Nil(t, r.result)
		require.NotNil(t, r.prev)
	})

	t.Run("miss 0", func(t *testing.T) {
		n := l.getTopNode()
		r := find(n, pointerTo(0), eqInt, gtInt)
		require.False(t, r.wasFound())
		require.Nil(t, r.result)
		require.NotNil(t, r.prev)
	})

	for _, value := range []int{4, 6, 9, 17, 22, 29} {
		t.Run(fmt.Sprintf("find %d", value), func(t *testing.T) {
			n := l.getTopNode()
			r := find(n, pointerTo(value), eqInt, gtInt)
			require.True(t, r.wasFound())
			require.NotNil(t, r.result)
			require.Equal(t, value, *r.result.bottomNode().value)
			require.NotNil(t, r.prev)
		})
	}
}
