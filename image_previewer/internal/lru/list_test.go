package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("complex2", func(t *testing.T) {
		l := NewList()

		for i, v := range [...]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80,60,40,20,0,10,30,50,70,90]

		require.Equal(t, 10, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 90, l.Back().Value)

		tmp := l.Front().Next.Next.Next.Next
		l.Remove(tmp)

		require.Equal(t, 9, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 90, l.Back().Value) // [80,60,40,20,10,30,50,70,90]
		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{80, 60, 40, 20, 10, 30, 50, 70, 90}, elems)

		for i := 0; i < l.Len(); i++ {
			l.MoveToFront(l.Back())
		}
		elems = make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{80, 60, 40, 20, 10, 30, 50, 70, 90}, elems)
	})

	t.Run("delete only one", func(t *testing.T) {
		l := NewList()

		el := l.PushFront(10)
		l.Remove(el)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}
