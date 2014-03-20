package main

import (
	"fmt"
	"github.com/gleicon/swarm"
)

var cluster = []string{"10.0.0.2"}
var myee = swarm.NewDistributedEventEmitter(cluster, "10.0.0.2")

func main() {
	myee.On("parrot", func(w []byte) {
		fmt.Println(w)
		myee.Emit("parrot", w)
	})
	for {
	}
}
