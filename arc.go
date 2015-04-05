package arc

import (
	"container/list"
)

type ARC struct {
	p     int
	c     int
	t1    *list.List
	b1    *list.List
	t2    *list.List
	b2    *list.List
	cache map[interface{}]*entry
}

func (a *ARC) req(ent *entry) {
	if ent.ll == a.t1 || ent.ll == a.t2 {
		// Case I
		ent.setMRU(a.t2)
	}
	if ent.ll == a.b1 {
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
	}
	if ent.ll == a.b2 {
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
	}
}

func (a *ARC) Put(key, value interface{}) bool {
	ent, ok := a.cache[key]
	if ok != true {
		// Case IV

		ent = &entry{
			key:   key,
			value: value,
		}

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
					a.replace(ent)
				}
			}
		}

		a.cache[key] = ent
		ent.setMRU(a.t1)
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
	return len(a.cache)
}

func (a *ARC) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	delete(a.cache, lru.Value.(*entry).key)
}

func (a *ARC) replace(ent *entry) {
	if a.t1.Len() > 0 && ((a.t1.Len() > a.p) || (ent.ll == a.b2 && a.t1.Len() == a.p)) {
		lru := a.t1.Back().Value.(entry)
		lru.setMRU(a.b1)
	} else {
		lru := a.t2.Back().Value.(entry)
		lru.setMRU(a.b2)
	}
}

func New(c int) *ARC {
	return &ARC{
		p:     0,
		c:     c,
		t1:    list.New(),
		b1:    list.New(),
		t2:    list.New(),
		b2:    list.New(),
		cache: make(map[interface{}]*entry, c),
	}
}
