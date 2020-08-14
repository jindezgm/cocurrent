/*
 * @Author: jinde.zgm
 * @Date: 2020-08-12 17:46:18
 * @Descripttion:
 */
package concurrent

import (
	"reflect"
	"sync"
	"testing"
)

type testMapUpdateValue struct {
	data []int
}

func TestMapUpdate(t *testing.T) {
	var m Map
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				m.Update("test", func(value interface{}) (interface{}, bool) {
					if nil == value {
						j--
						return testMapUpdateValue{data: make([]int, 10)}, true
					}
					old := value.(testMapUpdateValue)
					new := testMapUpdateValue{data: make([]int, 10)}
					copy(new.data, old.data)
					new.data[i]++
					return new, true
				})
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if value, exist := m.Load("test"); !exist {
		t.Fatal("value not exist")
	} else {
		v := value.(testMapUpdateValue)
		if len(v.data) != 10 {
			t.Fatal("invalid data len", len(v.data))
		}
		for i := range v.data {
			if v.data[i] != 100 {
				t.Fatal("invalid data value", i, v.data[i])
			}
		}
	}
}

type testRanger map[int]int

func (r *testRanger) Len() int { return len(*r) }
func (r *testRanger) Range(f func(key, value interface{}) bool) {
	for k, v := range *r {
		if !f(k, v) {
			break
		}
	}
}

func TestMapCopyAndClear(t *testing.T) {
	var m Map
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}

	r := make(testRanger)
	for i := 10; i < 20; i++ {
		r[i] = i
	}
	m.Copy(&r)
	for i := 0; i < 10; i++ {
		if value, exist := m.Load(i); exist {
			t.Fatal("old value not clear?", i, value)
		}
	}

	for k, v := range r {
		if value, exist := m.Load(k); !exist || !reflect.DeepEqual(v, value) {
			t.Fatal("value not copy?", k, exist, v, value)
		}
	}

	m.Clear()
	m.Range(func(key, value interface{}) bool {
		t.Fatal("any data after clear?", key, value)
		return false
	})
}
