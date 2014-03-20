package main

import (
	"fmt"
	"github.com/gleicon/swarm"
	"testing"
)

var cluster = []string{"10.0.0.1"}
var myee = NewDistributedEventEmitter(cluster, "10.0.0.1")

func main() {
	myee.On("parrot", func(w []byte) {
		fmt.Println(w)
		myee.Emit("parrot", w)
	})
}
