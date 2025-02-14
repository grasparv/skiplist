package skiplist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkiplistBasics(t *testing.T) {
	l := New(4, eqInt, gtInt)
	drawList(l)
	l.Insert(pointerTo(5))
	drawList(l)
	l.Insert(pointerTo(7))
	drawList(l)
	l.Insert(pointerTo(3))
	drawList(l)
	l.Remove(pointerTo(5))
	drawList(l)

	for i := 10; i < 100; i++ {
		l.Insert(pointerTo(i))
		require.True(t, l.Contains(pointerTo(i)))
		require.False(t, l.Contains(pointerTo(i+1)), i)
	}

	m := make(map[int]struct{})
	l.Range(pointerTo(1000), func(i int, value *int) bool {
		fmt.Printf("list %d\n", *value)
		_, ok := m[*value]
		require.False(t, ok)
		m[*value] = struct{}{}
		if i > 100 {
			return false
		}
		return true
	})
}

func TestSkiplistRandomOp(t *testing.T) {
	l := New(5, eqInt, gtInt)

	for i := 0; i < 1_000_000; i++ {
		x := rand.Intn(1_000)
		switch i % 4 {
		case 0:
			l.Insert(pointerTo(x))
		case 1:
			l.Remove(pointerTo(x))
		case 2:
			l.Contains(pointerTo(x))
		case 3:
			from := rand.Intn(1000)
			to := from + rand.Intn(100)
			l.Range(pointerTo(from), func(i int, value *int) bool {
				return i < to
			})
		}
	}
}

func TestSkiplistLarge(t *testing.T) {
	l := New(20, eqInt, gtInt)

	values := make(map[int]struct{})
	for i := 0; i < 5000; i++ {
		x := rand.Intn(1_000_000)
		values[x] = struct{}{}
		l.Insert(pointerTo(x))
	}

	require.Equal(t, len(values), l.Length())

	for i := 0; i < 1_000_000; i++ {
		_, ok := values[i]
		if !ok {
			require.False(t, l.Contains(pointerTo(i)))
		} else {
			require.True(t, l.Contains(pointerTo(i)))
		}
	}

	m := make(map[int]struct{})
	l.Range(pointerTo(1000), func(i int, value *int) bool {
		_, ok := m[*value]
		require.False(t, ok)
		m[*value] = struct{}{}
		if i > 100 {
			return false
		}
		return true
	})

	for v := range values {
		require.True(t, l.Contains(pointerTo(v)))
		l.Remove(pointerTo(v))
		require.False(t, l.Contains(pointerTo(v)))
	}
	require.Equal(t, 0, l.Length())

	for i := 0; i < 1_000_000; i++ {
		require.False(t, l.Contains(pointerTo(i)))
	}
}

func drawList(l *Skiplist[int]) {
	fmt.Printf("drawing..\n")
	for i := len(l.levels) - 1; i >= 0; i-- {
		for n := l.levels[i]; n != nil; n = n.next {
			if n.bottomNode().value == nil {
				fmt.Printf("H")
			} else {
				fmt.Printf("N")
			}
			if n.prev != nil {
				fmt.Printf("<")
			} else {
				fmt.Printf("-")
			}
			if n.down != nil {
				if n.bottomNode() != nil && n.bottomNode().value != nil {
					fmt.Printf("v%7d", *n.bottomNode().value)
				} else {
					fmt.Printf("v nil   ")
				}
			} else {
				if n.value != nil {
					fmt.Printf("X%7d", *n.value)
				} else {
					fmt.Printf("  base  ")
				}
			}
			if n.next != nil {
				fmt.Printf(">")
			} else {
				fmt.Printf("-")
			}
			fmt.Printf("    ")
		}
		fmt.Printf("\n")
	}
}
