AI Sandbox Go Bindings
======================

Go bindings to handle connection with http://aisandbox.com/ game server.

Instructions
------------

Download with

    go get github.com/errnoh/aisandbox


Then import the library to your code and call

    in, out, err := aisandbox.Connect("TerminatorKillerX", "serverhostname", port)

in: Channel where you receive parsed structs coming from the server.

out: Channel where you send your command structs

err: Possible error when connecting.

Notes
-----
* Since the command structs use anonymous structs to reduce the amount of structs we need, you can't create them with &struct syntax.

So instead of using

    &Attack{"Attack", {"Mr. Muggles", {2,3}, {15,26}, ""}}

use the old fashioned way

    attack := new(Attack) // or aisandbox.Attack
    attack.Value.Bot = "Pomerian"
    // etcetc..


* Non-fatal errors from the library are logged to standard logger)
* 'in' -channel will be closed when server sends <shutdown> message.
* Connection to the server will be closed from your end when you close 'out' -channel

TODO
----