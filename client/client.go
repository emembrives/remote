package main

import (
	"flag"
	"fmt"

	"github.com/emembrives/remote/proto"
	protobuf "github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

var (
	service  = flag.String("service", "", "Service name")
	endpoint = flag.String("endpoint", "", "Endpoint name")
	value    = flag.String("value", "", "Endpoint value")
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func send(requester *zmq.Socket, req *proto.Request) {
	data, err := protobuf.Marshal(req)
	panicOnErr(err)
	requester.SendBytes(data, 0)

	reply, err := requester.RecvBytes(0)
	panicOnErr(err)
	response := new(proto.Response)
	err = protobuf.Unmarshal(reply, response)
	panicOnErr(err)
	fmt.Println("Received ", response.String())
}

func main() {
	flag.Parse()

	requester, err := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	panicOnErr(err)
	requester.Connect("tcp://0.0.0.0:7001")

	req := &proto.Request{
		Type: proto.RequestType_SERVICE_DISCOVERY.Enum(),
	}
	send(requester, req)

	if len(*service) == 0 {
		return
	}

	req.Type = proto.RequestType_WRITE_ENDPOINT.Enum()
	req.WriteRequest = []*proto.Endpoint{&proto.Endpoint{Service: service,
		Endpoint: endpoint,
		Value:    value,
	}}
	send(requester, req)
}
