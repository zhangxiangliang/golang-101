# Web Crawler

## What is a web crawler?

* goal is to fetch all web pages, e.g. to feed to an indexer
* web pages and links form a graph
* multiple links to some pages
* graph has cycles

## Crawler challenges

### Exploit I/O concurrency
    
* Network latency is more limiting than network capacity
* Fetch many URLs at the same time, To increase URLs fetched per second
* Need threads for concurrency
  
### Fetch each URL only *once*

* avoid wasting network bandwidth
* be nice to remote servers
* Need to remember which URLs visited Know when finished

## Serial Crawler

* performs depth-first exploration via recursive Serial calls
* the "fetched" map avoids repeats, breaks cycles a single map, passed by reference, caller sees callee's updates
* but: fetches only one page at a time can we just put a "go" in front of the Serial() call? let's try it... what happened?

## ConcurrentMutex Crawler

* Creates a thread for each page fetch
* Many concurrent fetches, higher fetch rate
* the "go func" creates a goroutine and starts it running, func... is an "anonymous function"
* The threads share the "fetched" map, So only one thread will fetch any given page.

### Why the Mutex (Lock() and Unlock())?

#### One reason

* Two different web pages contain links to the same URL
* Two threads simultaneouly fetch those two pages
* T1 reads fetched[url], T2 reads fetched[url]
* Both see that url hasn't been fetched (already == false)
* Both fetch, which is wrong
* The lock causes the check and update to be atomic, So only one thread sees already==false.

#### Another reason:

* Internally, map is a complex data structure (tree? expandable hash?)
* Concurrent update/update may wreck internal invariants
* Concurrent update/read may crash the read

### What if I comment out Lock() / Unlock()?

* Why does it work `go run crawler.go` ? 
* Detects races even when output is correct `go run -race crawler.go`!

### How does the ConcurrentMutex crawler decide it is done?

* sync.WaitGroup Wait() waits for all Add()s to be balanced by Done()s i.e. waits for all child threads to finish
* there's a WaitGroup per node in the tree, How many concurrent threads might this crawler create?

## ConcurrentChannel Crawler

### A Go channel

* a channel is an object `ch := make(chan int)`
* a channel lets one thread send an object to another thread `ch <- x` the sender waits until some goroutine receives `y := <- ch`
* `for y := range ch` a receiver waits until some goroutine sends
* channels both communicate and synchronize
* several threads can send and receive on a channel
* channels are cheap
* remember: sender blocks until the receiver receives!
* "synchronous" watch out for deadlock

### ConcurrentChannel master()

* master() creates a worker goroutine to fetch each page
* worker() sends slice of page's URLs on a channel multiple workers send on the single channel
* master() reads URL slices from the channel

### At what line does the master wait?

* Does the master use CPU time while it waits?
* No need to lock the fetched map, because it isn't shared!

### How does the master know it is done?

* Keeps count of workers in n.
* Each worker sends exactly one item on channel.

### Why is it not a race that multiple threads use the same channel?

* Is there a race when worker thread writes into a slice of URLs, and master thread reads that slice, without locking?
* worker only writes slice *before* sending
* master only reads slice *after* receiving
* So they can't use the slice at the same time

### When to use sharing and locks, versus channels?

* Most problems can be solved in either style
* What makes the most sense depends on how the programmer thinks
    * state -- sharing and locks
    * communication -- channels
* For the 6.824 labs, I recommend sharing+locks for state, and sync.Cond or channels or time.Sleep() for waiting/notification.
