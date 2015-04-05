package main

import (
    "fmt"
    "container/list"
)

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func max(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

type ARC struct {
    p int
    c int
    T1 *list.List
    B1 *list.List
    T2 *list.List
    B2 *list.List
    cache map[interface{}]*entry
}

type entry struct {
	key interface{}
	value interface{}
	ll *list.List
	el *list.Element
}

func (e *entry) setLRU(list *list.List) {
	e.ll = list
	e.el = e.ll.PushBack(e)
}

func (e *entry) setMRU(list *list.List) {
	e.ll = list
	e.el = e.ll.PushFront(e)
}

func (e *entry) detach() {
	if e.ll != nil {
		e.ll.Remove(e.el)
	}
}

func (a *ARC) req(ent *entry) {
	if ent.ll == a.T1 || ent.ll == a.T2 {
		// Case I
		ent.setMRU(a.T2)
	}
	if ent.ll == a.B1 {
		// Case II
		// Cache Miss in T1 and T2
		
		// Adaptation
		var d int
        if a.B1.Len() >= a.B2.Len() {
            d = 1
        } else {
            d = a.B2.Len() / a.B1.Len()
        }
		a.p = min(a.p + d, a.c)

		a.replace(ent)
		ent.setMRU(a.T2)
	}
	if ent.ll == a.B2 {
		// Case III
		// Cache Miss in T1 and T2
		
		// Adaptation
		var d int
        if a.B2.Len() >= a.B1.Len() {
            d = 1
        } else {
            d = a.B1.Len() / a.B2.Len()
        }
		a.p = max(a.p - d, 0)

		a.replace(ent)
		ent.setMRU(a.T2)
	}
}

func (a *ARC) Put(key, value interface{}) bool {
	ent, ok := a.cache[key]
	if ok != true {
		// Case IV

		ent = &entry{
			key: key,
			value: value,
		}

		if a.T1.Len() + a.B1.Len() == a.c {
			// Case A
			if a.T1.Len() < a.c {
				a.delLRU(a.B1)
				a.replace(ent)
			} else {
				a.delLRU(a.T1)
			}
		} else if a.T1.Len() + a.B1.Len() < a.c {
			// Case B
			if a.T1.Len() + a.T2.Len() + a.B1.Len() + a.B2.Len() >= a.c {
				if a.T1.Len() + a.T2.Len() + a.B1.Len() + a.B2.Len() == 2*a.c {
					a.delLRU(a.B2)
					a.replace(ent)
				}
			}
		}

		a.cache[key] = ent
		ent.setMRU(a.T1)
	} else {
		ent.value = value
		a.req(ent)
	}
	return ok
}

func (a *ARC) Get(key interface{}) (value interface{}, ok bool) {
	ent, ok := a.cache[key]
	if ok {
		a.req(ent)
		return ent.value, true
	}
	return nil, ok
}

func (a *ARC) Len() int {
	return a.T1.Len() + a.T2.Len() + a.B1.Len() + a.B2.Len()
}

func (a *ARC) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	delete(a.cache, lru.Value.(*entry).key)
}

func (a *ARC) replace(ent *entry) {
	if a.T1.Len() > 0 && ((a.T1.Len() > a.p) || (ent.ll == a.B2 && a.T1.Len() == a.p)) {
		lru := a.T1.Back().Value.(entry)
		lru.setMRU(a.B1)
	} else {
		lru := a.T2.Back().Value.(entry)
		lru.setMRU(a.B2)
	}
}

func New(c int) *ARC {
	return &ARC{
		p: 0,
		c: c,
		T1: list.New(),
		B1: list.New(),
		T2: list.New(),
		B2: list.New(),
		cache: make(map[interface{}]*entry, c),
	}
}

func main() {

	cache := New(3)

	cache.Put("Hello", "World")
	cache.Put("Hello2", "World2")
	cache.Put("Hello3", "World3")
	cache.Put("Hello", "World")


	fmt.Println(cache.Get("Hello"))

    fmt.Println("Hello World")
}
