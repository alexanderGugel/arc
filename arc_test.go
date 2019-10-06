package arc

import "testing"

func TestBasic(t *testing.T) {
	cache := New(3)
	if got, want := cache.Len(), 0; got != want {
		t.Errorf("empty cache.Len(): got %d want %d", cache.Len(), want)
	}

	const (
		k1 = "Hello"
		k2 = "Hallo"
		k3 = "Ciao"
		k4 = "Salut"

		v1 = "World"
		v2 = "Worlds"
		v3 = "Welt"
	)

	// Insert the first value
	cache.Put(k1, v1)
	if got, want := cache.Len(), 1; got != want {
		t.Errorf("insertion of key #%d: cache.Len(): got %d want %d", want, cache.Len(), want)
	}
	if got, ok := cache.Get(k1); !ok || got != v1 {
		t.Errorf("cache.Get(%q): got (%q, %t) want (%q, true)", k1, got, ok, v1)
	}

	// Replace existing value for a given key
	cache.Put(k1, v2)
	if got, want := cache.Len(), 1; got != want {
		t.Errorf("re-insertion: cache.Len(): got %d want %d", cache.Len(), want)
	}
	if got, ok := cache.Get(k1); !ok || got != v2 {
		t.Errorf("re-insertion: cache.Get(%q): got (%q, %t) want (%q, true)", k1, got, ok, v2)
	}

	// Add a second different key
	cache.Put(k2, v3)
	if got, want := cache.Len(), 2; got != want {
		t.Errorf("insertion of key #%d: cache.Len(): got %d want %d", want, cache.Len(), want)
	}
	if got, ok := cache.Get(k1); !ok || got != v2 {
		t.Errorf("cache.Get(%q): got (%q, %t) want (%q, true)", k1, got, ok, v2)
	}
	if got, ok := cache.Get(k2); !ok || got != v3 {
		t.Errorf("cache.Get(%q): got (%q, %t) want (%q, true)", k2, got, ok, v3)
	}

	// Fill cache
	cache.Put(k3, v1)
	if got, want := cache.Len(), 3; got != want {
		t.Errorf("insertion of key #%d: cache.Len(): got %d want %d", want, cache.Len(), want)
	}

	// Exceed size, this should evict:
	cache.Put(k4, v1)
	if got, want := cache.Len(), 3; got != want {
		t.Errorf("insertion of key out of size: cache.Len(): got %d want %d", cache.Len(), want)
	}
}
