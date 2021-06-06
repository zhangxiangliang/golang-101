
# Key-value

* A toy key/value storage server -- Put(key,value), Get(key)->value
* Uses Go's RPC library

## Common

* Declare Args and Reply struct for each server handler.

## Client

* connect()'s Dial() creates a TCP connection to the server get() and put() are client "stubs"
* call() asks the RPC library to perform the call
* you specify server function name, arguments, place to put reply
* library marshalls args, sends request, waits, unmarshalls reply
* return value from call() indicates whether it got a reply
* usually you'll also have a reply.Err indicating service-level failure

## Server

* Go requires server to declare an object with methods as RPC handlers
* Server then registers that object with the RPC library
* Server accepts TCP connections, gives them to RPC library
* The RPC library
    * reads each request
    * creates a new goroutine for this request
    * unmarshalls request
    * looks up the named object (in table create by register())
    * calls the object's named method (dispatch)
    * marshalls reply
    * writes reply on TCP connection
* The server's Get() and Put() handlers
    * Must lock, since RPC library creates a new goroutine for each request
    * read args; modify reply

## A few details

### Binding

* how does client know what server computer to talk to?
* For Go's RPC, server name/port is an argument to Dial
* Big systems have some kind of name or configuration server
  
### Marshalling

* format data into packets
* Go's RPC library can pass strings, arrays, objects, maps, &c
* Go passes pointers by copying the pointed-to data
* Cannot pass channels or functions
