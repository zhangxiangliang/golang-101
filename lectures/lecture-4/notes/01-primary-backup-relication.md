# Primary/Backup Replication

* Primary/Backup Replication for Fault Tolerance, case study of VMware FT, an extreme version of the idea.
* The topic is still fault tolerance to provide availability despite server and network failures using replication.

## What kinds of failures can replication deal with?

### "fail-stop" failure of a single replica
    
* fan stops working, CPU overheats and shuts itself down someone trips over replica's power cord or network cable software
* notices it is out of disk space and stops

### Maybe not defects in h/w or bugs in s/w or human configuration errors

* Often not fail-stop
* May be correlated (i.e. cause all replicas to crash at the same time).
* But sometimes can be detected (e.g. checksums)

### How about earthquake or city-wide power failure?

* Only if replicas are physically separated

## Two main replication approaches:

### State transfer

* Primary replica executes the service
* Primary sends new state to backups

### Replicated state machine

* Clients send operations to primary, primary sequences and sends to backups
* All replicas execute all operations
* If same start state, same operations, same order, deterministic, then same end state.

## Difference In Two main replication approaches.

## State transfer is simpler

* But state may be large, slow to transfer over network.

## Replicated state machine

* Operations are often small compared to state
* But complex to get right
* VM-FT uses replicated state machine

## Big Questions

* What state to replicate?
* Does primary have to wait for backup?
* When to cut over to backup?
* Are anomalies visible at cut-over?
* How to bring a replacement backup up to speed?

## At what level do we want replicas to be identical?

### Application state, e.g. a database's tables?

* GFS works this way
* Can be efficient; primary only sends high-level operations to backup.
* Application code server must understand fault tolerance, to e.g forward op stream

### Machine level, e.g. registers and RAM content?

* might allow us to replicate any exising server w/o modification
* requires forwording of machine events (interrupts, DMA, &c)
* requires machine modifications to send/recv event stream

