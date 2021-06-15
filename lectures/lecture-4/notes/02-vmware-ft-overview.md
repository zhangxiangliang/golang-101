# VMware FT

* Transparent: can run any existing OS and server software
* Appears like a single server to clients

## Overview

* diagram: app, OS, VM-FT underneath, disk server, network, clients
* hypervisor == monitor == VMM (virtual machine monitor)
* OS and App is the guest running inside a virtual machine
* two machines , primary and backup.
* primary sends all external events (client packets) to backup over network "logging channel", carrying log entries
* ordinarily, backup's output is suppressed by FT
* if either stops being able to talk to the other over the network
    * "goes live" and provides sole service
    * if primary goes live, it stops sending log entries to the backup

## VMM emulates a local disk interface

* but actual storage is on a network server.
* treated much like a client:
    * usually only primary communicates with disk server (backup's FT discards)
    * if backup goes lvie, it talks to disk server
* external disk makes creating a new backup faster(don't have to copy primary's disk)

## When does the primary have to send information to the backup?

* Any time someting happens that might cause their executions to diverge.
* Anything that's not a deterministic consequence  of executing instructions.

## What sources of divergence must FT handle?

* Must instructions excute indentically on primary and backup.
* As long as memory registers are identical,  which we're assuming by induction.
* Inputs from external world -- just network packets. These appear as DMA'd data plus an interrupt.
* Timing of interrupts.
* Instructions that aren't functions of state, such as reading current time.
* Not multi-core races, since uniprocessor only.

## Why would divergence be a disaster?

b/c state on backup would differ from state on primary, and if primary then failed, clients would see inconsistency. Example: GFS lease expiration: 

* Imagine we're replicating the GFS master
* Chunkserver must send "please renew" msg before 60-second lease expires
* Clock interrupt drives master's notion of time
* Suppose chunkserver sends "please renew" just around 60 seconds
* On primary, clock interrupt happens just after request arrives. Primary copy of master renews the lease, to the same chunkserver.
* On backup, clock interrupt happens just before request. Backup copy of master expires the lease.
* If primary fails, backup takes over, it will think there is no lease, and grant it to a different chunkserver. Then two chunkservers will have lease for same chunk.

So backup must see same events, in same order, at same points in instruction stream.