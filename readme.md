AI Sandbox Go Bindings
======================

Go bindings to handle connection with http://aisandbox.com/ game server.

Instructions
------------

Download with

    go get github.com/errnoh/aisandbox


Then import the library to your code and call Connect() with host, port and bot name as parameters.

    in, out, err := aisandbox.Connect("serverhostname", port, "TerminatorKillerX")

in: Channel where you receive parsed structs coming from the server.

out: Channel where you send your Commands

err: Possible error when connecting.

Notes
-----
* Non-fatal errors from the library are logged to standard logger)
* 'in' -channel will be closed when server sends <shutdown> message.
* Connection to the server will be closed from your end when you close 'out' -channel
* in is type <-chan interface (receive only)
* out is type chan<- aisandbox.Command (send only)
* there are predefined structs for each of the four commands that satisfy aisandbox.Command interface.

TODO
----
