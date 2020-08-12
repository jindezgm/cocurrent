package concurrent

import (
	"testing"

	"github.com/google/uuid"
)

func TestMap(t *testing.T) {
	var m Map

	key := uuid.New().String()
	if ok := m.Update(key, func(value interface{}) (interface{}, bool) {
		if nil != value {
			t.Fatal("first update value not nil")
		}
		return "first", true
	}); !ok {
		t.Fatal("update failed")
	}
	if value, ok := m.Load(key); !ok || value.(string) != "first" {
		t.Fatal("update failed", value)
	}

	if ok := m.Update(key, func(value interface{}) (interface{}, bool) {
		if "first" != value.(string) {
			t.Fatal("second update value not first", value)
		}
		return "second", true
	}); !ok {
		t.Fatal("update failed")
	}
	if value, ok := m.Load(key); !ok || value.(string) != "second" {
		t.Fatal("update failed", value)
	}

	if ok := m.Update(key, func(value interface{}) (interface{}, bool) {
		if "second" != value.(string) {
			t.Fatal("third update value not second", value)
		}
		return "third", false
	}); ok {
		t.Fatal("update failed")
	}
	if value, ok := m.Load(key); !ok || value.(string) != "second" {
		t.Fatal("update failed", value)
	}

	m.Delete(key)
	if ok := m.Update(key, func(value interface{}) (interface{}, bool) {
		if nil != value {
			t.Fatal("fourth update value not nil", value)
		}
		return "fourth", true
	}); !ok {
		t.Fatal("update failed")
	}
	if value, ok := m.Load(key); !ok || value.(string) != "fourth" {
		t.Fatal("update failed", value)
	}
}
