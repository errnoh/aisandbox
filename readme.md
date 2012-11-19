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


Notes
-----

* Easy way to wrap the commands is to use function signature that looks something like:

    func attack(name string, direction []float64, description string, coords ...[]float64)

That way you can call the command with as many waypoints as you can. See the bot in _example for.. example.

* _example folder contains a sample bot which you can use as an example.
* Non-fatal errors from the library are logged to standard logger)
* 'in' -channel will be closed when server sends <shutdown> message.
* Connection to the server will be closed from your end when you close 'out' -channel
* in is type <-chan interface (receive only)
* out is type chan<- aisandbox.Command (send only)
* there are predefined structs for each of the four commands that satisfy aisandbox.Command interface.

Support
-------

During the AI Sandbox CTF competition, real time support can be found (but not promised) from the official IRC channel.
I should be reachable from 1200 to 2200 GMT almost every day. Just ask a question, highlight or open a query.

For actual issues considering the bindings feel free to open an issue on Github, or again, just notify me on IRC.

Bindings are usually updated within couple hours of official game updates. Give or take couple hours.

TODO
----
