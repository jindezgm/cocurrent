/*
 * @Author: jinde.zgm
 * @Date: 2020-08-12 17:46:18
 * @Descripttion:
 */
package concurrent

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type testMapUpdateValue struct {
	data []int
}

type revisionedValue struct {
	Revision int64
	Value    interface{}
}

func TestMapUpdate(t *testing.T) {
	parallel := runtime.NumCPU()
	runtime.GOMAXPROCS(parallel)
	var m Map
	var wg sync.WaitGroup
	wg.Add(parallel)
	var count int32
	for i := 0; i < parallel; i++ {
		go func(i int) {
			for j := i; j < 100; j++ {
				if m.Update("test", func(value interface{}) (interface{}, bool) {
					if nil == value {
						return revisionedValue{Revision: 1}, true
					}
					rv := value.(revisionedValue)
					if rv.Revision == 100 {
						return revisionedValue{}, false
					}
					return revisionedValue{Revision: rv.Revision + 1}, true
				}) {
					t.Log("update ok", i, atomic.AddInt32(&count, 1))
				} else {
					break
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if count != 100 {
		t.Fatal(count)
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
