package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

//func (lruCache *lruCache) Set(key Key, value interface  {}){
//	lruCache.items[key] = &ListItem{
//		Value: value,
//		Next:  nil,
//		Prev:  nil,
//	}
//	lruCache.queue.PushBack(lruCache.items[key])
//}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
