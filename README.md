# Bifrost

#### Bifrost is a transport-agnostic message broker engine

## Goals

The main goal of this project is to provide a robust messaging solution that is
simple to maintain, yet still scalable to (at least) medium size deployments.
To this end the project attempts to provide sane defaults when ever possible.

Currently, the Engine is only capable of running on a single node, but plans
are in the works to utilize a Raft algorithm to maintain a distributed queue
for multi-node deployments.

## Using The Client

Bifrost includes a native Go client library for interacting with the HTTP
portions of the API.

```go
import (
	"fmt"
	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	name := "test.channel"
	size := 10

	// First, create client object and do any error checking
	cl, _ := bifrost.New("127.0.0.1", 2727)


	// Second, create a new channel on the server to send messages through.
	resp := cl.Do(bifrost.CreateChannel(name, size))
	if err := resp.Error(); err != nil {
		// handle error...
	}

	// Now that the channel has been created, we can see it in the API
	var channel bifrost.Channel
	resp = cl.Do(bifrost.GetChannel(name))
	if err := resp.Decode(&channel); err != nil {
		// handle error...
	}
	fmt.Println("Name=%v\nMax=%v", channnel.Name, channel.Max)

	// Let's try publishing a message.
	resp = cl.Do(
		channel.Publish("Very important message to pass")
	)
	if err := resp.Error(); err != nil {
		// handle error...
	}

	// Finally, pop the message off of the queue.
	var message []byte
	err = cl.Do(channel.Pop).Decode(&message)
	if err != nil {
		// handle error...
	}
	fmt.Println("Message =", string(message))
```

Godocs for the client can be found here:
https://godoc.org/github.com/goslang/bifrost/client

More examples can be found in the examples directory, here:
https://github.com/goslang/bifrost/tree/master/examples

## Using The Engine

See the Godocs for in depth library documentation:
https://godoc.org/github.com/goslang/bifrost/engine

## Building the HTTP/Websocket server

Just run `make`:
```bash
make
./bifrost-server
```
