package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	lock     sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, isKeyExists := l.items[key]

	if isKeyExists {
		l.queue.Remove(l.items[key])
	} else if l.queue.Len() >= l.capacity {
		cacheItem, ok := l.queue.Back().Value.(cacheItem)
		if !ok {
			return false
		}

		l.queue.Remove(l.items[cacheItem.key])
		delete(l.items, cacheItem.key)
	}

	l.items[key] = l.queue.PushFront(cacheItem{key: key, value: value})

	return isKeyExists
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	listItem, ok := l.items[key]

	if !ok {
		return nil, false
	}

	cacheItem, ok := listItem.Value.(cacheItem)
	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(listItem)

	return cacheItem.value, true
}

func (l *lruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
