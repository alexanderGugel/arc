// Package arc implements an Adaptive replacement cache
package arc

import (
	"container/list"
	"sync"
)

type ARC struct {
	p     int
	c     int
	t1    *list.List
	b1    *list.List
	t2    *list.List
	b2    *list.List
	mutex sync.RWMutex
	len   int
	cache map[interface{}]*entry
}

// New returns a new Adaptive Replacement Cache (ARC).
func New(c int) *ARC {
	return &ARC{
		p:     0,
		c:     c,
		t1:    list.New(),
		b1:    list.New(),
		t2:    list.New(),
		b2:    list.New(),
		len:   0,
		cache: make(map[interface{}]*entry, c),
	}
}

// Put inserts a new key-value pair into the cache.
// This optimizes future access to this entry (side effect).
func (a *ARC) Put(key, value interface{}) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok != true {
		a.len++

		ent = &entry{
			key:   key,
			value: value,
			ghost: false,
		}

		a.req(ent)
		a.cache[key] = ent
	} else {
		if ent.ghost {
			a.len++
		}
		ent.value = value
		ent.ghost = false
		a.req(ent)
	}
	return ok
}

// Get retrieves a previously via Set inserted entry.
// This optimizes future access to this entry (side effect).
func (a *ARC) Get(key interface{}) (value interface{}, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok {
		a.req(ent)
		return ent.value, !ent.ghost
	}
	return nil, false
}

// Len determines the number of currently cached entries.
// This method is side-effect free in the sense that it does not attempt to optimize random cache access.
func (a *ARC) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.len
}

func (a *ARC) req(ent *entry) {
	if ent.ll == a.t1 || ent.ll == a.t2 {
		// Case I
		ent.setMRU(a.t2)
	} else if ent.ll == a.b1 {
		// Case II
		// Cache Miss in t1 and t2

		// Adaptation
		var d int
		if a.b1.Len() >= a.b2.Len() {
			d = 1
		} else {
			d = a.b2.Len() / a.b1.Len()
		}
		a.p = min(a.p+d, a.c)

		a.replace(ent)
		ent.setMRU(a.t2)
	} else if ent.ll == a.b2 {
		// Case III
		// Cache Miss in t1 and t2

		// Adaptation
		var d int
		if a.b2.Len() >= a.b1.Len() {
			d = 1
		} else {
			d = a.b1.Len() / a.b2.Len()
		}
		a.p = max(a.p-d, 0)

		a.replace(ent)
		ent.setMRU(a.t2)
	} else if ent.ll == nil {
		// Case IV

		if a.t1.Len()+a.b1.Len() == a.c {
			// Case A
			if a.t1.Len() < a.c {
				a.delLRU(a.b1)
				a.replace(ent)
			} else {
				a.delLRU(a.t1)
			}
		} else if a.t1.Len()+a.b1.Len() < a.c {
			// Case B
			if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() >= a.c {
				if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() == 2*a.c {
					a.delLRU(a.b2)
				}
				a.replace(ent)
			}
		}

		ent.setMRU(a.t1)
	}
}

func (a *ARC) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	a.len--
	delete(a.cache, lru.Value.(*entry).key)
}

func (a *ARC) replace(ent *entry) {
	if a.t1.Len() > 0 && ((a.t1.Len() > a.p) || (ent.ll == a.b2 && a.t1.Len() == a.p)) {
		lru := a.t1.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b1)
	} else {
		lru := a.t2.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
	}
}
