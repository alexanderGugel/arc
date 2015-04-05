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

func TestBasicReplace(t *testing.T) {
	cache := New(3)

	cache.Put("Hello", "Hallo")
	cache.Put("World", "Welt")
	cache.Get("World")
	cache.Put("Cache", "Cache")
	cache.Put("Replace", "Ersetzen")

	value, ok := cache.Get("World")
	if !ok || value != "Welt" {
		t.Error("ARC should have replaced \"Hello\"")
	}

	if cache.Len() != 3 {
		t.Error("ARC should have a maximum size of 3")
	}
}
