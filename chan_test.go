/*
 * @Author: jinde.zgm
 * @Date: 2020-08-20 20:50:06
 * @Descripttion:
 */
package concurrent

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestQueuedChan(t *testing.T) {
	parallel := runtime.NumCPU()
	runtime.GOMAXPROCS(parallel)

	c := NewQueuedChan()
	wg := sync.WaitGroup{}
	wg.Add(parallel * 10)
	for i := 0; i < parallel; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				c.PushChan() <- i*10 + j
			}
		}(i)
		go func(i int) {
			for {
				x, ok := <-c.PopChan()
				if !ok {
					return
				}
				t.Log(i, x)
				wg.Done()
			}
		}(i)
	}
	wg.Wait()
	if c.Len() != 0 {
		t.Fatal("test queued channel failed", c.Len())
	}
	c.Close()

	c = NewQueuedChan()
	for i := 0; i < 10; i++ {
		c.Push(i)
	}
	time.Sleep(time.Millisecond * 10)
	if count := c.Remove(func(i interface{}) (bool, bool) {
		if i.(int)%2 == 0 {
			return true, true
		}
		return false, true
	}); count != 5 {
		t.Fatal("remove failed", count)
	} else if c.Len() != 5 {
		t.Fatal("remove failed", c.Len())
	}

	for i := 0; i < c.Len(); i++ {
		x := c.Pop().(int)
		if x%2 == 0 {
			t.Fatal("remove failed", x)
		}
	}
}
