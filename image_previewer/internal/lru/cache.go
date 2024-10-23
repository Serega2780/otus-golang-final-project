package lru

import "sync"

type Key string

type Pair struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	mux      *sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

func NewPair(k Key, v interface{}) *Pair {
	return &Pair{
		key:   k,
		value: v,
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		mux:      &sync.Mutex{},
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Clear() {
	lru.mux.Lock()
	defer lru.mux.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mux.Lock()
	defer lru.mux.Unlock()
	val, ok := lru.items[key]
	if ok {
		item, _ := val.Value.(Pair)
		item.value = value
		val.Value = item
		lru.queue.MoveToFront(val)
	} else {
		if lru.queue.Len() == lru.capacity {
			back := lru.queue.Back()
			lru.queue.Remove(back)
			v := back.Value.(Pair).key
			delete(lru.items, v)
		}
		lru.items[key] = lru.queue.PushFront(*NewPair(key, value))
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mux.Lock()
	defer lru.mux.Unlock()
	val, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(val)
		return val.Value.(Pair).value, ok
	}
	return nil, ok
}
