package main

import (
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

	req := &proto.Request{}
	data, err := protobuf.Marshal(req)
	requester.SendBytes(data, 0)
}
