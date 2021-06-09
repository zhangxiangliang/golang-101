# Each log entry: instruction #, type, data.

## FT's handling of timer interrupts

### Goal

* primary and backup should see interrupt at the same point in the instruction stream
  
### Primary

* FT fields the timer interrupt
* FT reads instruction number from CPU
* FT sends "timer interrupt at instruction X" on logging channel
* FT delivers interrupt to primary, and resumes it (this relies on CPU support to interrupt after the X'th instruction)

### Backup

* ignores its own timer hardware
* FT sees log entry *before* backup gets to instruction X
* FT tells CPU to interrupt (to FT) at instruction X
* FT mimics a timer interrupt to backup

## FT's handling of network packet arrival (input)

### Primary

* FT tells NIC to copy packet data into FT's private "bounce buffer"
* At some point NIC does DMA, then interrupts
* FT gets the interrupt
* FT pauses the primary
* FT copies the bounce buffer into the primary's memory
* FT simulates a NIC interrupt in primary
* FT sends the packet data and the instruction # to the backup

### Backup

 * FT gets data and instruction # from log stream
 * FT tells CPU to interrupt (to FT) at instruction X
 * FT copies the data to backup memory, simulates NIC interrupt in backup

## Why the bounce buffer?

* We want the data to appear in memory at exactly the same point in execution of the primary and backup.
* Otherwise they may diverge.

## Note that the backup must lag by one one log entry

Suppose primary gets an interrupt, or input, after instruction X:

* If backup has already executed past X, it cannot handle the input correctly
* So backup FT can't start executing at all until it sees the first log entry
* Then it executes just to the instruction # in that log entry And waits for the next log entry before resuming backup.

### Example: non-deterministic instructions

Some instructions yield different results even if primary/backup have same state,e.g. reading the current time or cycle count or processor serial #

#### Primary:

* FT sets up the CPU to interrupt if primary executes such an instruction
* FT executes the instruction and records the result
* sends result and instruction # to backup

#### Backup:

* FT reads log entry, sets up for interrupt at instruction #
* FT then supplies value that the primary got

## What about output (sending network packets)?

* Primary and backup both execute instructions for output
* Primary's FT actually does the output
* Backup's FT discards the output

### Output example: DB server

clients can send "increment" request DB increments stored value, replies with new value
  
so:
* suppose the server's value starts out at 10
* network delivers client request to FT on primary
* primary's FT sends on logging channel to backup
* FTs deliver request to primary and backup
* primary executes, sets value to 11, sends "11" reply, FT really sends reply
* backup executes, sets value to 11, sends "11" reply, and FT discards
* the client gets one "11" response, as expected

But wait:

* suppose primary crashes just after sending the reply so client got the "11" reply
* AND the logging channel discards the log entry w/ client request primary is dead, so it won't re-send
* backup goes live but it has value "10" in its memory!
* now a client sends another increment request it will get "11" again, not "12" oops

Solution: the Output Rule (Section 2.2)
  
* before primary sends output, must wait for backup to acknowledge all previous log entries

Again, with output rule primary:

* receives client "increment" request
* sends client request on logging channel
* about to send "11" reply to client
* first waits for backup to acknowledge previous log entry
* then sends "11" reply to client
* suppose the primary crashes at some point in this sequence
* if before primary receives acknowledgement from backup,maybe backup didn't see client's request, and didn't increment, but also primary won't have replied
* if after primary receives acknowledgement from backup then client may see "11" reply, but backup guaranteed to have received log entry w/ client's request, so backup will increment to 11.

The Output Rule is a big deal

* Occurs in some form in all replication systems
* A serious constraint on performance
* An area for application-specific cleverness
* Eg. maybe no need for primary to wait before replying to read-only operation
* FT has no application-level knowledge, must be conservative

## What if the primary crashes just after getting ACK from backup, but before the primary emits the output? Does this mean that the output won't ever be generated?

* Here's what happens when the primary fails and the backup goes live.
* The backup got some log entries from the primary.
* The backup continues executing those log entries WITH OUTPUT DISCARDED.
* After the last log entry, the backup goes live -- stops discarding output
* In our example, the last log entry is arrival of client request
* So after client request arrives, the client will start emitting outputs
* And thus it will emit the reply to the client

## But what if the primary crashed *after* emitting the output? Will the backup emit the output a *second* time?

* OK for TCP, since receivers ignore duplicate sequence numbers.
* OK for writes to disk, since backup will write same data to same block #.
* Duplicate output at cut-over is pretty common in replication systems
* Clients need to keep enough state to ignore duplicates
* Or be designed so that duplicates are harmless

## Does FT cope with network partition -- could it suffer from split brain? E.g. if primary and backup both think the other is down. Will they both go live?

* The disk server breaks the tie.
    * Disk server supports atomic test-and-set.
    * If primary or backup thinks other is dead, attempts test-and-set.
    * If only one is alive, it will win test-and-set and go live.
    * If both try, one will lose, and halt.

* The disk server may be a single point of failure
    * If disk server is down, service is down
    * They probably have in mind a replicated disk server

## Why don't they support multi-core?

* FT/Non-FT: impressive! little slow down
* Logging bandwidth, Directly reflects disk read rate + network input rate 18 Mbit/s for my-sql
* These numbers seem low to me, Applications can read a disk at at least 400 megabits/second, So their applications aren't very disk-intensive

## When might FT be attractive?
  
* Critical but low-intensity services, e.g. name server.
* Services whose software is not convenient to modify.

## What about replication for high-throughput services?

* People use application-level replicated state machines for e.g. databases.
    * The state is just the DB, not all of memory+disk.
    * The events are DB commands (put or get), not packets and interrupts.
* Result: less fine-grained synchronization, less overhead.
* GFS use application-level replication, as do Lab 2 &c

## Summary:

* Primary-backup replication VM-FT: clean example
* How to cope with partition without single point of failure?

## Next lecture

* How to get better performance?
* Application-level replicated state machines