package arc

import "testing"

func TestBasic(t *testing.T) {
    cache := New(3)
    if cache.Len() != 0 {
        t.Error("Empty cache should have length 0")
    }

    cache.Put("Hello", "World")
    if cache.Len() != 1 {
        t.Error("Cache should have length 1")
    }

    var val interface{}
    var ok bool

    if val, ok = cache.Get("Hello"); val != "World" || ok != true {
        t.Error("Didn't set \"Hello\" to \"World\"")
    }

    cache.Put("Hello", "World1")
    if cache.Len() != 1 {
        t.Error("Inserting the same entry multiple times shouldn't increase cache size")
    }

    if val, ok = cache.Get("Hello"); val != "World1" || ok != true {
        t.Error("Didn't update \"Hello\" to \"World1\"")
    }
    
    cache.Put("Hallo", "Welt")
    if cache.Len() != 2 {
        t.Error("Inserting two different entries should result into lenght=2")
    }

    if val, ok = cache.Get("Hallo"); val != "Welt" || ok != true {
        t.Error("Didn't set \"Hallo\" to \"Welt\"")
    }
}
