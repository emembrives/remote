package main

import (
	"fmt"

	"github.com/emembrives/tinkerings/android-remote/proto"
	protobuf "github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	requester, err := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	panicOnErr(err)
	requester.Connect("tcp://0.0.0.0:7001")

	req := &proto.Request{
		Type: proto.RequestType_SERVICE_DISCOVERY.Enum(),
	}
	data, err := protobuf.Marshal(req)
	panicOnErr(err)
	requester.SendBytes(data, 0)

	reply, err := requester.RecvBytes(0)
	panicOnErr(err)
	fmt.Println("Received ", reply)
}
