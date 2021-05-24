package lru

import (
	"fmt"
	"testing"
)

type person struct {
	name string
	age  int
}

func (p *person) Len() int {
	return int(len(p.name) + len(fmt.Sprint("v", p.age)))
}

var person1 = &person{"ZhangSan", 28}
var person2 = &person{"LisiLisiLisiLisi", 100}

func TestGet(t *testing.T) {

	cache := NewCache(4<<10, nil)
	t.Run("test add and get", func(t *testing.T) {
		cache.Add(person1.name, person1)

		value, ok := cache.Get(person1.name)
		if !ok {
			t.Fatal("failed get the value ")
		}

		if value.Len() != person1.Len() {
			t.Fatalf("want %v value but got %v", person1, value)
		}
	})

	t.Run("testing update and get", func(t *testing.T) {
		cache.Add(person1.name, person1)
		cache.Add(person1.name, person2)

		value, ok := cache.Get(person1.name)
		if !ok {
			t.Fatal("failed get the value ")
		}

		if value.Len() == person1.Len() {
			t.Fatalf("falied update of add method")
		}
	})
}

func TestRemoveOldest(t *testing.T) {
	cache := NewCache(4<<10, nil)

	cache.Add(person1.name, person1)
	cache.Add(person2.name, person2)

	cache.RemoveOldest()

	if v, ok := cache.Get(person1.name); ok {
		t.Fatalf("remove old value failed, get %v", v)
	}
	if v, ok := cache.Get(person2.name); !ok {
		t.Fatalf("remove old value is wrong, new value [%v] is deleted", v)
	}
}
