package golru

import (
	"container/list"
	"errors"
	"sync"
)

// Cache is a structure describing a temporary data store by key through a linked list. Mutex is present inside to
// implement the ability to work with the cache competitively.
// The number of entries is fixed. This implementation, rather than a memory limit, is chosen to work directly with
// storage objects. In addition to the list, there is a hash table with data inside the structure, which is used
// directly for accounting records.
// All fields are non-exportable, which allows you to work with the content through methods without having
// direct access to the cache
type Cache struct {
	mu sync.Mutex

	capacity int
	items    map[string]*list.Element
	chain    *list.List
}

// NewCache create new implementation of lru Cache. Capacity can't be less than one. If you set capacity to zero,
// for example, an assignment error will return
func NewCache(n int) (*Cache, error) {
	if n <= 0 {
		return nil, errors.New("capacity of the cache can not be less than 1")
	}
	return &Cache{
		capacity: n,
		items:    make(map[string]*list.Element, n),
		chain:    list.New(),
	}, nil
}

// item is an element inside *list.Element of Cache with the key and value used by the your program
type item struct {
	key   string
	value interface{}
}
