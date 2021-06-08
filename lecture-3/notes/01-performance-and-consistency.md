# The Google File System

## Why are we reading this paper?

* distributed storage is a key abstraction
    * what should the interface/semantics look like?
    * how should it work internally?
* GFS paper touches on many themes of 6.824
    * parallel performance
    * fault tolerance
    * replication
    * consistency
* good systems paper -- details from app all the way to network
* successful real-world design

## Why is distributed storage hard?

* high performance -> shard data over many servers
* many servers -> constant faults
* fault tolerance -> replication
* replication -> potential inconsistencies
* better consistency -> low performance

## What would we like for consistency?

Ideal model: same behavior as a signle server:

* server uses disk storage
* server executes client operations one at time (even if concurrent)
* reads reflect previous writes even if server crashes and restarts
* thus: suppose C1 and C2 write concurrently, and after the writes have completed, C3 and C4 read. what can they see?

| client | options |
| --- | --- |
| C1 | Write x equal 1 |
| C2 | Write x equal 2 |
| C3 | Read x |
| C4 | Read x |

* Either 1 or 2, but both have to see the same value. This is a "strong" consistency model.
* But a single server has poor fault-tolerance.

## Replication for fault-tolerance makes strong consistency tricky.

A simple but broken replication scheme:

* Two replica servers, S1 and S2
* clients send writes to both, in parallel
* clients send reads to either

In our example, C1 and C2 write messages could arrive in different orders at the two replicas

* if C3 reads S1, it might see x = 1
* if C4 reads S2, it might see x = 2

* or what if S1 receives a write, but the client crashes before sending the write to S2? 
* That's not strong consistency!
* Better consistency usually requires communication to ensure the replicas stay in sync -- can be slow!
* Lots of tradeoffs possible between performance and consistency.