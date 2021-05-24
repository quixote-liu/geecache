package geecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()

	if g, ok := groups[name]; ok {
		return g
	}
	return nil
}

func (g *Group) Get(key string) (ByteViews, error) {
	if key == "" {
		return ByteViews{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Geecache hit]")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteViews, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (value ByteViews, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteViews{}, err
	}

	value = ByteViews{b: cloneBytes(bytes)}
	g.populataCache(key, value)
	return value, nil
}

func (g *Group) populataCache(key string, value ByteViews) {
	g.mainCache.add(key, value)
}
