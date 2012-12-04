AI Sandbox Go Bindings
======================

Go bindings to handle connection with http://aisandbox.com/ game server.

Instructions
------------

Download with

    go get github.com/errnoh/aisandbox

Then import the library to your code with

    import "github.com/errnoh/aisandbox"

 and call Connect() with host, port and bot name as parameters.

    in, out, err := aisandbox.Connect("serverhostname", port, "TerminatorKillerX")

in: Channel where you receive parsed structs coming from the server.

out: Channel where you send your Commands

err: Possible error when connecting.

Server then sends you one aisandbox.LevelInfo and one aisandbox.GameInfo struct to process before starting the game.
When you're done processing those, use:

    aisandbox.Ready()

To inform the server that you're ready to start the game.

NOTE FOR THE COMPETITION:
Your bot should support custom host and port as command-line arguments.
Hostname as first argument and port number as second.
See example bot for.. example.

Constructors
------------

Since JSON API 1.2, there are constructors for each Command type.
Each one accepts any amount of []float64 coordinates as last parameter.
Coordinates are in []float64{x,y} format
NewDefend also allows third value inside the slice, resulting in slice that looks like []float64{x,y,duration}

See example bot for sample usage.
(Example doesn't pass multiple waypoints for move/charge/attack commands but those are supported as well.)

Notes
-----

* _example folder contains a sample bot which you can use as an example.
* Non-fatal errors from the library are logged to standard logger)
* 'in' -channel will be closed when server sends <shutdown> message.
* Connection to the server will be closed from your end when you close 'out' -channel
* in is type <-chan interface (receive only)
* out is type chan<- aisandbox.Command (send only)

Support
-------

During the AI Sandbox CTF competition, real time support can be found (but not promised) from the official IRC channel.
I should be reachable from 1200 to 2200 GMT almost every day. Just ask a question, highlight or open a query.

For actual issues considering the bindings feel free to open an issue on Github, or again, just notify me on IRC.

Bindings are usually updated within couple hours of official game updates. Give or take couple hours.

TODO
----
