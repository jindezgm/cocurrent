<!--
 * @Author: jinde.zgm
 * @Date: 2020-08-12 21:59:48
 * @Descripttion: 
-->
# Introduce
Concurrent is a package designed for concurrency.
## Map
Map extend Update interface based on sync.Map, equivalent to CAS(Compare And Swap).  
Update keeps calling 'tryUpdate()' to update key 'key' retrying the update until success if there is conflict. Note that value passed to tryUpdate may change across invocations of tryUpdate() if other writers are simultaneously updating it, so tryUpdate() needs to take into account the current contents of the value when deciding how the update object should look. If the key doesn't exist, it will pass nil to tryUpdate.
```go
    // Map example
    var m concurrent.Map
    if ok := m.Update("test", func(value interface{}){
        if nil == value {
            return 1, true
        }
        if i := value.(int); i != 1 {
            return nil, false
        }
        return 2, true
    }); ok {
        fmt.Println("update ok")
    }
```
## QueuedChan
QueuedChan is a channel with dynamic buffer size. Unlick traditional channel, push is never blocked and the buffer size follows the performance changes of production and consumption.
```go
    // QueuedChan example
    c := concurrent.NewQueuedChan()
    defer c.Close()
    select{
    case c.PushChan() <- 0:
    }
    c.Pop()
    fmt.Println(c.Len())
    c.Remove(func(i interface) (bool, bool) {
        if i.(int) == 0 {
            return true, false
        }
        return false, true
    })
```