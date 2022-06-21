package LRU

import "container/list"

type Cache struct {
	maxBytes  int64 //	允许使用的最大内存
	nBytes    int64 //	当前已使用的最大内存
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) //	记录删除(缓存中不存在)时的回调函数
}

type Entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

//	Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 	查询
func (c *Cache) Get(key string) (value Value, ok bool) {
	if elm, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elm)
		kv := elm.Value.(*Entry)
		return kv.value, true
	}
	return
}

// 移除LRU节点
func (c *Cache) RemoveOldest() {
	elm := c.ll.Back()
	if elm != nil {
		c.ll.Remove(elm)
		kv := elm.Value.(*Entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
func (c *Cache) Add(key string, value Value) {
	if elm, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elm)
		kv := elm.Value.(*Entry)
		c.nBytes += int64(len(key)) - int64(kv.value.Len())
		kv.value = value
	} else {
		elm := c.ll.PushFront(&Entry{key, value})
		c.cache[key] = elm
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// number of list member
func (c *Cache) Len() int {
	return c.ll.Len()
}
