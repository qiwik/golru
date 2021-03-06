package golru

import (
	"container/list"
	"reflect"
)

// Add returns false if current key already exists, and true if key doesn't exist and new item was added to cache.
// When a new element is added, it is placed at the top of the list, and if capacity is reached, the last element,
// which is also the most unpopular in the cache, is deleted
func (c *Cache) Add(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.validate(key); ok {
		return false
	}

	if c.chain.Len() == c.capacity {
		c.removeLast()
	}

	newItem := &item{
		key:   key,
		value: value,
	}
	newElement := c.chain.PushFront(newItem)
	c.items[newItem.key] = newElement
	return true
}

// Get func returns a value with true if such element exist with current key, else returns nil and false. If an element
// exists, it is moved to the top of the list in the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return nil, false
	}

	value := element.Value.(*item).value
	c.chain.MoveToFront(element)
	return value, true
}

// Remove returns false if current key doesn't exist, and true if removing from cache was successful
func (c *Cache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return false
	}

	delete(c.items, key)
	c.chain.Remove(element)
	return true
}

// ChangeValue allows you to change the value of a key that already exists in the cache. If there is no such key in
// the cache, the function returns false. If the value has changed, the element is sent to the top of the cache list
func (c *Cache) ChangeValue(key string, newValue interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return false
	}

	element.Value.(*item).value = newValue
	c.chain.MoveToFront(element)
	return true
}

// validate checks the existence of an element by the key, and if it does not exist, returns false, instead of an element
func (c *Cache) validate(key string) (element *list.Element, ok bool) {
	if element, ok = c.items[key]; !ok {
		return nil, false
	}
	return element, true
}

// Clear completely clears the cache
func (c *Cache) Clear() {
	for c.chain.Len() > 0 {
		c.removeLast()
	}
}

// Len allows you to find out the fullness of the cache
func (c *Cache) Len() int {
	return c.chain.Len()
}

// ChangeCapacity allows you to dynamically change the cache capacity. The new value must not be less than one. If
// the new capacity is less than the previous one, then the last elements in the list are deleted up to the desired
// parameter value
func (c *Cache) ChangeCapacity(newCap int) {
	switch {
	case newCap <= 0:
		return
	case newCap >= c.capacity:
		c.capacity = newCap
		return
	default:
		c.capacity = newCap
		for c.Len() > newCap {
			c.removeLast()
		}
	}
}

// removeLast deletes the last element in the list
func (c *Cache) removeLast() {
	currentElement := c.chain.Back()
	last := c.chain.Remove(currentElement).(*item)
	delete(c.items, last.key)
}

// Keys returns a slice of the keys that exist in the cache by simply traversing all the keys. Works faster than
// a function with reflection
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// ReflectKeys returns a slice of keys existing in the cache using reflection. It works 3-4 times slower than the Keys
// function, but is left for variability
func (c *Cache) ReflectKeys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	keysValues := reflect.ValueOf(c.items).MapKeys()
	keys := make([]string, 0, len(c.items))
	for i := range keysValues {
		keys = append(keys, keysValues[i].String())
	}
	return keys
}

// Values returns a slice of all existing element values in the cache
func (c *Cache) Values() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]interface{}, 0, len(c.items))
	for _, value := range c.items {
		values = append(values, value)
	}
	return values
}
