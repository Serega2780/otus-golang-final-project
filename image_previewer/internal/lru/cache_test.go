package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache, _ := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache, _ = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache, _ = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
		c := NewCache(3)

		wasInCache, _ := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache, _ = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache, _ = c.Set("ccc", 300)
		require.False(t, wasInCache)
		wasInCache, oldest := c.Set("ddd", 400)
		require.False(t, wasInCache)
		require.Equal(t, 100, oldest.(int))

		_, wasInCache = c.Get("aaa")
		require.False(t, wasInCache)
		_, wasInCache = c.Get("bbb")
		require.True(t, wasInCache)
		_, wasInCache = c.Get("ccc")
		require.True(t, wasInCache)
		_, wasInCache = c.Get("ddd")
		require.True(t, wasInCache)
	})

	t.Run("purge old logic", func(t *testing.T) {
		// Write me
		c := NewCache(3)

		wasInCache, _ := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache, _ = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache, _ = c.Set("ccc", 300)
		require.False(t, wasInCache)

		c.Get("aaa")
		c.Get("bbb")
		c.Get("ccc") // order in a queue [ccc, bbb, aaa]

		c.Set("aaa", 101) // order in a queue [aaa, ccc, bbb]
		c.Set("ccc", 101) // order in a queue [ccc, aaa, bbb]

		c.Set("ddd", 400) // order in a queue [ddd, aaa, ccc]
		_, wasInCache = c.Get("aaa")
		require.True(t, wasInCache)
		_, wasInCache = c.Get("bbb")
		require.False(t, wasInCache)
		_, wasInCache = c.Get("ccc")
		require.True(t, wasInCache)
		_, wasInCache = c.Get("ddd")
		require.True(t, wasInCache)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
