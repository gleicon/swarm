package swarm

import (
	"fmt"
	"testing"
)

var cluster = []string{"10.0.0.1"}
var myee = NewDistributedEventEmitter(cluster, "10.0.0.1")

func TestDistributedEventEmitter(t *testing.T) {
	myee.On("jazz", func(wot []byte) {
		if string(wot) != "el test" {
			t.Error("Invalid value")
		}
	})
	myee.Emit("jazz", []byte("el test"))
}

func TestDistributedEventEmitterRemoveAllisteners(t *testing.T) {
	myee.RemoveAllListeners("jazz")
	l := myee.Listeners("jazz")
	if len(l) != 0 {
		t.Error("Invalid listeners value")
	}
}

func TestDistributedEventEmitterListeners(t *testing.T) {
	myee.RemoveAllListeners("jazz")
	myee.On("jazz", func(wot []byte) { fmt.Println(wot) })
	myee.On("jazz", func(wot []byte) { fmt.Println(wot) })
	l := myee.Listeners("jazz")
	if len(l) != 2 {
		t.Error("Invalid listeners value")
	}
}

func TestRemoveAllListeners(t *testing.T) {
	myee.On("jazz", func(wot []byte) { fmt.Println(wot) })
	myee.On("jazz", func(wot []byte) { fmt.Println(wot) })
	myee.RemoveAllListeners("jazz")
	if len(myee.Listeners("jazz")) != 0 {
		t.Error("Invalid listeners value")
	}
	myee.On("jazz", func(wot []byte) {
		if string(wot) != "el test" {
			t.Error("Invalid value")
		}
	})
	myee.Emit("jazz", []byte("el test"))
}

func BenchmarkDistributedEventEmitterSequential(b *testing.B) {
	myee.On("jazz", func(wot []byte) { fmt.Println(string(wot)) })
	for i := 0; i < b.N; i++ {
		st := fmt.Sprintf("iteration: %d", i)
		myee.Emit("jazz", []byte(st))
	}
}

func BenchmarkDistributedEventEmitterGoroutines(b *testing.B) {
	myee.On("jazz", func(wot []byte) { fmt.Println(string(wot)) })
	for i := 0; i < b.N; i++ {
		go func(idx int) {
			st := fmt.Sprintf("goroutine: %d", idx)
			myee.Emit("jazz", []byte(st))
		}(i)
	}
}

func TestCluster(t *testing.T) {
	myee2 := NewDistributedEventEmitter(cluster, "10.0.0.2")
	n := myee2.Nodes()
	if len(n) != 2 {
		t.Error("Invalid node value", len(n))
	}
}
