package lru

import (
	"fmt"
	"testing"
)

type person struct {
	name string
	age  int
}

func (p *person) Len() int64 {
	return int64(len(p.name) + len(fmt.Sprint("v", p.age)))
}

func TestGet(t *testing.T) {
	person1 := &person{"ZhangSan", 28}
	person2 := &person{"LisiLisiLisiLisi", 100}

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
