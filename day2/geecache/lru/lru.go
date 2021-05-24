package lru

import "container/list"

type Cache struct {
	maxBytes int64

	nbytes int64

	ll *list.List

	cache map[string]*list.Element

	// 回调函数
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)

		// 判断字典中的值的类型
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	// 选择最后一个
	ele := c.ll.Back()

	// 在字典中进行删除
	if ele != nil {
		c.ll.Remove(ele)

		// 在字典中删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)

		// 更新内存大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 修改
		c.ll.MoveToFront(ele)

		kv := ele.Value.(*entry)

		c.nbytes += int64(value.Len()) - int64(kv.value.Len())

		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})

		c.cache[key] = ele

		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	if c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}
