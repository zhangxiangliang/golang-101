# Remote Procedure Call (RPC)

* a key piece of distributed system machinery; all the labs use RPC
* goal: easy-to-program client/server communication
* hide details of network protocols
* convert data (strings, arrays, maps, &c) to "wire format"

## RPC message diagram:

```
Client ----request----> Server
Client <---response---- Server
```

## Software structure

```
client app        handler fns
stub fns          dispatcher
RPC lib           RPC lib
net  ------------ net
```

## RPC problem

### What to do about failures?

e.g. lost packet, broken network, slow server, crashed server

### What does a failure look like to the client RPC library?

* Client never sees a response from the server
* Client does *not* know if the server saw the request!
* Maybe server never saw the request
* Maybe server executed, crashed just before sending reply
* Maybe server executed, but network died just before delivering reply

### Simplest failure-handling scheme: "best effort"
  
* Call() waits for response for a while
* If none arrives, re-send the request
* Do this a few times
* Then give up and return an error

#### Is "best effort" easy for applications to cope with?

* A particularly bad situation: client executes `Put("k", 10); Put("k", 20);` both succeed. What will Get("k") yield? [diagram, timeout, re-send, original arrives late]

#### Is best effort ever OK?

* read-only operations, operations that do nothing if repeated, e.g. DB checks if record has already been inserted

### Better RPC behavior "At most once" 

server RPC code detects duplicate requests returns previous reply instead of re-running handler.

#### How to detect a duplicate request?

client includes unique ID (XID) with each request uses same XID for re-send

```
if seen[xid]:
    r = old[xid]
else
    r = handler()
    old[xid] = r
    seen[xid] = true
```

#### What if two clients use the same XID?

* big random number?
* combine unique client ID (ip address?) with sequence #?

### Server must eventually discard info about old RPCs when is discard safe?

* each client has a unique ID (perhaps a big random number)
* per-client RPC sequence numbers
* client includes "seen all replies <= X" with every RPC
* much like TCP sequence #s and acks
* or only allow client one outstanding RPC at a time
* arrival of seq+1 allows server to discard all <= se

### How to handle dup req while original is still executing?

* server doesn't know reply yet
* "pending" flag per executing RPC; wait or ignore

### What if an at-most-once server crashes and re-starts?

* if at-most-once duplicate info in memory, server will forget and accept duplicate requests after re-start.
* maybe it should write the duplicate info to disk
* maybe replica server should also replicate duplicate info

### Go RPC is a simple form of "at-most-once" ?

* open TCP connection
* write request to TCP connection
* Go RPC never re-sends a request So server won't see duplicate requests
* Go RPC code returns an error if it doesn't get a reply
    * perhaps after a timeout (from TCP)
    * perhaps server didn't see request
    * perhaps server processed request but server/net failed before reply came back

### What about "exactly once"?

* unbounded retries plus duplicate detection plus fault-tolerant service