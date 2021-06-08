# GFS

## Context

* Many Google services needed a big fast unified storage system MapReduce, crawler/indexer, log storage/analysis, Youtobe.
* Global (over a single data center): any client can read any file 
* Allows sharding of data among applications
* Automatic "sharding" of each file over many server/disks: For parallel performance, To increase space available.
* Automatic recovery from failures
* Just one data center per deployment
* Just Google aplications/users
* Aimed at sequential access to huge files, read or append. Not a low-latency DB for small items.

## What was new about this in 2003? How did they get an SOSP paper accepted?

* Not the basic ideas of distribution, sharding, fault-tolerance
* Huge scale.
* Used in industry, real-world experience
* Successful use of weak consistency
* Successful use of single master

