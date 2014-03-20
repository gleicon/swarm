### Swarm

Swarm is a distributed event emitter written in Go. It uses memberlist to orchestrate a cluster of eventemitters and zeromq to exchange serialized messages between the nodes. Inside each node it provides a familiar EventEmitter interface to attach functions to events and emit data.

### Tests and benchmark

    $ go test -v
    $ go test -bench=.

### Usage


