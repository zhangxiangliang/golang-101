# Overall structure

* clients (library, RPC -- but not visible as a UNIX File System)
* each file split into independent 64 MB chunks
* chunk servers, each chunk replicated on 3
* every file's chunks are spread over the chunk servers for parallel read/write (e.g. MapReduce), and to allow huge files
* single master, and master replicas.
* division of work: master deals w/ naming, chunk servers w/ data.

## Mastar state In RAM (for speed, must be smallish):

* filename -> array of chunk handles (Non-Volatile)

* chunk handle
    * version # (Non-Volatile)
    * list of chunk servers (Volatile)
    * primary (Volatile)
    * lease time (Volatile)

* on disk
    * log
    * checkpoint

* questions:
    * Why a log?
    * Why a checkpoint?
    * Why big chunks?

## What are the steps when client C wants to read a file?

* C sends filename and offset to master M (if not cached)
* M finds chunk handle for that offset
* M replies with list of chunk servers only those with lastest version
* C caches hanle and chunk server list
* C sends request to nearest chunk server, chunk handle, offset.
* chunk server reads from chunk file on disk, returns.

## What are the steps when C wants to do a "record append"?

* C asks M about file's last chunk
* if M sees chunk ha no primary (or lease expired)
    * if no chunk servers w/ lastest version #, error
    * pick primary and secondaries from those w/ lastest version #
    * increment version #, write to log on disk
    * tell primary and secondaries who they are, anad new version #
    * replicas write new version # to disk
* M tells C the primary and secondaries
* C sends data to all (just temporary), waits
* C tells M to append
* M checks that lease hasn't expired, and chunk has space
* M picks an offset (at end of chunk)
* M writes chunk file (a Linux file)
* M tells each secondary the offset, tells to append to chunk file
* M waits for all secondaries to reply, or timeout secondary can reply "error" e.g. out of disk space
* M tells C "ok" or "error"
* C retries from start if error

## What consistency guarantees does GFS provide to clients?

Needs to be in a form that tells applications how to use GFS. Here's a possibility:

* If the primary tells client that a record append succeeded, then any reader that subsequently openss the file and scans it will sedd the appended record some where.
* But not that failed appends won't be visible, or that all readers will see the same file content, or the same order of records.

## Summary
  case study of performance, fault-tolerance, consistency, specialized for MapReduce applications

### good ideas:
    * global cluster file system as universal infrastructure
    * separation of naming (master) from storage (chunkserver)
    * sharding for parallel throughput
    * huge files/chunks to reduce overheads
    * primary to sequence writes
    * leases to prevent split-brain chunkserver primaries

### not so great:
    * single master performance
    * ran out of RAM and CPU
    * chunk servers not very efficient for small files
    * lack of automatic fail-over to master replica
    * maybe consistency was too relaxed