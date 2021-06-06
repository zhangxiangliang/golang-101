# Infrastructure: RPC and threads

Threads and RPC in Go, with an eye towards the labs.

## Why Go?

* Good support for threads
* Convenient RPC
* Type safe
* Garbage-collected (no use after freeing problems)
* threads and GC is particularly attractive!
* relatively simple
* After the tutorial, use http://golang.org/doc/effective_go.html

## Threads

* A useful structuring tool, but can be tricky
* Go calls them goroutines; everyone else calls them threads

## Threads of execution

* threads allow one program to do many things at once
* each thread executes serially, just like an ordinary non-threaded program
* the threads share memory
* each thread includes some per-thread state: program counter, registers, stack

## Why Threads

They express concurrency, which you need in distributed systems

### I/O Concurrency

* Client sends requests to many severs in parallel and waits for replies.
* Server processes multiple client request; each request many block.
* While waiting for the disk to read data for client x, process a request from client y.

### Multicore performance

* Execute code in parallel on several cores.

### Convenience

In background, once per second, check whether each worker is still alive.


## Is there an alternative to threads?

* Write code that explicitly interleaves activities, in a single thread. Usually called "Event-driven".
* Keep a table of state about each activity, e.g. each client request.
* One "event" loop that: checks for new input for each activity (e.g. arrival of reply from server), does the next step for each activity, updates state.
* Event-driven gets you I/O concurrency, and eliminates threads costs (which can be substantial), but doesn't get multi-core speedup, and is painful to program.

## Threading challenges:

### Shared data

e.g. what if two threads do n = n + 1 at the same time? or one thread reads while another increments? 

#### Question 

This is a "RACE" and usually a bug

#### Solution 

* use locks (Go's sync.Mutex)
* or avoid sharing mutable data

### Coordination between threads

e.g. one thread is producing data, another thread is consuming it:

#### Question

* How can the consumer wait (and release the CPU) ?
* How can the producer wake up the consumer?

#### Solution

* use Go channels or sync.Cond or WaitGroup

### Deadlock

cycles via locks and/or communication (e.g. RPC or Go channels)