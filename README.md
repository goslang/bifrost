# Bifrost

#### Bifrost is a transport-agnostic message broker engine

## Goals

The main goal of this project is to provide a robust messaging solution that is
simple to maintain, yet still scalable to (at least) medium size deployments.
To this end the project attempts to provide sane defaults when ever possible.

Currently, the Engine is only capable of running on a single node, but plans
are in the works to utilize a Raft algorithm to maintain a distributed queue
for multi-node deployments.

## Using The Engine

See the Godocs for in depth library documentation:
https://godoc.org/github.com/goslang/bifrost/engine

## Building the HTTP/Websocket server

Just run `make`:
```bash
make
./bifrost-server
```
