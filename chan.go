/*
 * @Author: jinde.zgm
 * @Date: 2020-08-20 20:11:32
 * @Descripttion:
 */
package concurrent

import (
	"container/list"
	"sync/atomic"
)

// QueuedChan is a channel with dynamic buffer size implemented using a list
type QueuedChan struct {
	*list.List                  // Object list
	len        int32            // Because List.Len() is not thread safe, atomic access to the variable to implement thread safety
	close      chan struct{}    // Close signal
	done       chan struct{}    // Closed signal
	pushc      chan interface{} // Push object channel
	popc       chan interface{} // Pop object channel
	ctrlc      chan interface{} // Control channel
}

// NewQueuedChan create QueuedChan.
func NewQueuedChan() *QueuedChan {
	c := &QueuedChan{
		List:  list.New(),
		close: make(chan struct{}),
		done:  make(chan struct{}),
		pushc: make(chan interface{}),
		popc:  make(chan interface{}),
		ctrlc: make(chan interface{}),
	}
	// Create a coroutine to push and pop object.
	go c.run()

	return c
}

// queuedChanRemoveCmd is remove object command.
type queuedChanRemoveCmd struct {
	f func(i interface{}) (ok bool, cont bool) // Object filter function.
	r chan int                                 // Object number removed.
}

// PushChan get push object channel.
func (c *QueuedChan) PushChan() chan<- interface{} { return c.pushc }

// Push object into channel.
func (c *QueuedChan) Push(i interface{}) {
	select {
	case c.pushc <- i:
	case <-c.close:
	}
}

// PopChan get pop object channel.
func (c *QueuedChan) PopChan() <-chan interface{} { return c.popc }

// Pop object from channel.
func (c *QueuedChan) Pop() interface{} {
	select {
	case i := <-c.popc:
		return i
	case <-c.close:
		return nil
	}
}

// Len get buffered object count.
func (c *QueuedChan) Len() int { return int(atomic.LoadInt32(&c.len)) }

// Remove the object filtered by function f.
// The two return values of function f: the first is whether to delete, the second is whether to continue.
func (c *QueuedChan) Remove(f func(i interface{}) (ok bool, cont bool)) int {
	// Send remove command.
	r := make(chan int)
	select {
	case c.ctrlc <- queuedChanRemoveCmd{f: f, r: r}:
	case <-c.close:
		return 0
	}

	// Waiting for result.
	select {
	case n := <-r:
		return n
	case <-c.close:
		return 0
	}
}

// Close channel
func (c *QueuedChan) Close() {
	close(c.close)
	<-c.done

	// Close pop channel avoid reader blocked.
	close(c.popc)
}

// CloseAndFlush close channel and flush buffered objects.
func (c *QueuedChan) CloseAndFlush() {
	close(c.close)
	<-c.done

	// Flush buffered objects.
	c.flush()

	// Close pop channel avoid reader blocked.
	close(c.popc)
}

// run get object from push channel and push back into list,
// pop front object from list and output by pop channel.
func (c *QueuedChan) run() {
	// Notify close channel coroutine.
	defer close(c.done)

	for {
		var elem *list.Element
		var item interface{}
		var popc chan<- interface{}

		// Get front element of the queue.
		if elem = c.Front(); nil != elem {
			popc, item = c.popc, elem.Value
		}

		select {
		// Put the new object into the end of queue.
		case i := <-c.pushc:
			c.PushBack(i)
		// Remove the front element from queue if send out success
		case popc <- item:
			c.List.Remove(elem)
		// Control command
		case cmd := <-c.ctrlc:
			c.control(cmd)
		// Channel is closed
		case <-c.close:
			return
		}
		// Update channel length
		atomic.StoreInt32(&c.len, int32(c.List.Len()))
	}
}

// flush flush the value buffered in the queue.
func (c *QueuedChan) flush() {
	// Flush queue.
	for elem := c.Front(); nil != elem; elem = c.Front() {
		c.popc <- elem.Value
		c.List.Remove(elem)
	}
}

// control channel.
func (c *QueuedChan) control(ctrl interface{}) {
	switch cmd := ctrl.(type) {
	// Remove object from channel.
	case queuedChanRemoveCmd:
		c.remove(&cmd)
	// Unknown command.
	default:
		panic(ctrl)
	}
}

// remove execute remove command
func (c *QueuedChan) remove(cmd *queuedChanRemoveCmd) {
	// Get object count before remove.
	count := c.List.Len()
	// Iterate list.
	for i := c.Front(); i != nil; {
		var re *list.Element
		// Filter object.
		ok, cont := cmd.f(i.Value)
		if ok {
			re = i
		}
		// Next element.
		i = i.Next()
		// Remove element
		if nil != re {
			c.List.Remove(re)
		}
		// Continue
		if !cont {
			break
		}
	}
	// Update channel length
	atomic.StoreInt32(&c.len, int32(c.List.Len()))
	// Return removed object number.
	cmd.r <- count - c.List.Len()
}
