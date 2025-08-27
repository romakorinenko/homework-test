package hw04lrucache

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

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("1", 1)
		cache.Set("2", 2)
		cache.Set("3", 3)

		val, ok := cache.Get("1")
		require.True(t, ok)
		require.Equal(t, 1, val)
		val, ok = cache.Get("2")
		require.True(t, ok)
		require.Equal(t, 2, val)
		val, ok = cache.Get("3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		cache.Set("4", 4)
		val, ok = cache.Get("1")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = cache.Get("2")
		require.True(t, ok)
		require.NotNil(t, val)
		val, ok = cache.Get("3")
		require.True(t, ok)
		require.NotNil(t, val)
		val, ok = cache.Get("4")
		require.True(t, ok)
		require.NotNil(t, val)
	})

	t.Run("purge with changes", func(t *testing.T) {
		cache := NewCache(3)
		cache.Set("1", 1)
		cache.Set("2", 2)
		cache.Set("3", 3)

		cache.Set("1", 10)
		cache.Set("2", 20)
		cache.Set("3", 30)

		val, ok := cache.Get("1")
		require.True(t, ok)
		require.Equal(t, 10, val)
		val, ok = cache.Get("2")
		require.True(t, ok)
		require.Equal(t, 20, val)
		val, ok = cache.Get("3")
		require.True(t, ok)
		require.Equal(t, 30, val)

		cache.Set("4", 4)
		val, ok = cache.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		cache.Set("5", 5)
		cache.Set("6", 6)
		val, ok = cache.Get("2")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = cache.Get("3")
		require.False(t, ok)
		require.Nil(t, val)

		cache.Clear()
		val, ok = cache.Get("4")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
