package swarm

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/ugorji/go/codec"
	"net"
	"os"
)

const (
	maxDGSize = 1024
	mcastAddr = "224.0.0.251:9000"
)

type Document map[string]interface{}

type DistributedEventEmitter struct {
	nodeId    string
	listeners map[string][]func([]byte)
	ml        *memberlist.Memberlist
	sub       *net.UDPConn
}

func NewDistributedEventEmitter(cluster []string, bindAddr string) *DistributedEventEmitter {
	dee := DistributedEventEmitter{}
	c := memberlist.DefaultLANConfig()
	c.Name = bindAddr
	c.BindAddr = bindAddr

	ml, err := memberlist.Create(c)
	dee.ml = ml

	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	_, err = dee.ml.Join(cluster)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
	h, err := os.Hostname()
	if err != nil {
		panic("Failed to get hostname" + err.Error())
	}

	fmt.Sprintf(dee.nodeId, "%s:%d", h, os.Getpid())

	dee.listeners = make(map[string][]func([]byte))
	a, err := net.ResolveUDPAddr("udp", mcastAddr)

	if err != nil {
		panic("Error converting mcast addr: " + err.Error())
	}

	dee.sub, err = net.ListenMulticastUDP("udp", nil, a)
	dee.sub.SetReadBuffer(maxDGSize)

	if err != nil {
		panic("Failed listen to UDP mcast: " + err.Error())
	}

	go dee.readLoop(dee.sub)

	return &dee
}

func (dee DistributedEventEmitter) readLoop(s *net.UDPConn) {
	buf := make([]byte, maxDGSize)
	for {
		rlen, remote, err := s.ReadFromUDP(buf)
		fmt.Println("buffer: ", buf)
		fmt.Println("remote addr: ", remote)
		fmt.Println("rlen: ", rlen)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		d := make(Document)
		var mh codec.MsgpackHandle
		dec := codec.NewDecoderBytes(buf, &mh)
		err = dec.Decode(&d)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		if d["nodeId"] != dee.nodeId {
			fmt.Println(d)
		}
	}
}

func (dee DistributedEventEmitter) sendUDPData(wot []byte) {
	addr, err := net.ResolveUDPAddr("udp", mcastAddr)
	if err != nil {
		fmt.Println("Error resolvind mcast addr: ", err.Error())
	}
	c, err := net.DialUDP("udp", nil, addr)
	d := make(Document)
	d["nodeId"] = dee.nodeId
	d["payload"] = wot

	b := []byte{}
	var mh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&b, &mh)
	err = enc.Encode(d)
	c.Write(b)
}

func (dee DistributedEventEmitter) On(eventname string, f func(payload []byte)) {
	dee.listeners[eventname] = append(dee.listeners[eventname], f)
}

func (dee DistributedEventEmitter) RemoveAllListeners(eventname string) {
	dee.listeners[eventname] = nil
}

func (dee DistributedEventEmitter) Listeners(eventname string) []func([]byte) {
	return dee.listeners[eventname]
}

func (dee DistributedEventEmitter) Nodes() []*memberlist.Node {
	return dee.ml.Members()
}

func (dee DistributedEventEmitter) Emit(eventname string, message []byte) {
	for f := range dee.listeners[eventname] {
		fun := dee.listeners[eventname][f]
		if fun != nil {
			fun(message)
		}
		dee.sendUDPData(message)
	}
	for _, member := range dee.ml.Members() {
		fmt.Printf("Sending message to: %s %s\n", member.Name, member.Addr)
	}
}
