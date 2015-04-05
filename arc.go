package main

import (
    "fmt"
    "github.com/alexanderGugel/golang-lru"
)

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

type ARC struct {
    p int
    c int
    T1 lru.Cache
    B1 lru.Cache
    T2 lru.Cache
    B2 lru.Cache
}

func (c *ARC) Add(key, value interface{}) {
    if c.T1.Has(key) || c.T2.Has(key) {
        c.T1.Remove(key)
        c.T2.Add(key, value)
    } else if c.B1.Has(key) {
        if c.B1.Len() >= c.B2.Len() {
            d := 1
        } else {
            d := c.B2.Len() / c.B1.Len()
        }
        c.p = min(c.p + d, c.c)

        c.replace(key, value, c.p)
        c.B1.Remove(key)
        c.T2.Add(key, value)
    } else if c.B2.Has(key) {
        if c.B2.Len() >= c.B1.Len() {
            d := 1
        } else {
            d := c.B1.Len() / c.B2.Len()
        }
        c.p = min(c.p - d, 0)
        // TODO replace
    } else {
        if c.T1.Len() + c.B1.Len() == c.c {
            if c.T1.Len() < c.c {
                // TODO
            } else {
                // TODO
            }
        } else if c.T1.Len() + c.B1.Len() < c.c {
        	if c.T1.Len() + c.T2.Len() + c.B2.Len() + c.B1.Len() >= c.c {
        		// TODO
        	}
        }
        // TODO
    }
}

// Get looks up a key's value from the cache.
func (c *ARC) Get(key interface{}) (value interface{}, ok bool) {
	// TODO
}

func (c *ARC) replace(key interface{}, value interface{}, p int) {
	if c.T1.Len() > 0 && (c.T1.Len() > c.p || (c.B2.Has(key) && c.T1.Len() == c.p)) {
		c.T1.Remove(key)
		c.B1.Add(key, value)
	} else {
		c.T2.Remove(key)
		c.B2.Add(key, value)
	}
}


func main() {
    fmt.Println("Hello World")
}
