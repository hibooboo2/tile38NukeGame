# A basic game with boxes moving around in 2d plane.
Support Native UI as well as browser ui.


go build should compile everything to all be in one thing.

There are 2 parts the server.. Which serves a js ui that is the same ui as the native ui just compiled as GOOS=js and GOARCH=wasm

All assests should be compiled into binary for server and clients.
Server includes the client that is for web.. EG it will have a pre compiled wasm client that can be served.
Both native ui and js ui will need a server to talk to.

Native ui And js ui should abstract out creation of ui elements and events so that both can use the same code for the game logic itself.
This will be done using a wrapper gui framework that we create. It will have methods for Drawing elements we want etc.. Or us render style loop for things.. Will depend on performance etc...

Should support keyboard events and mouse events.
