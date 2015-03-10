package main

import (
	"fmt"
	"log"
	"os"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/emembrives/tinkerings/android-remote/proto"

	zmq "github.com/pebbe/zmq4"
)

const (
	url = "tcp://127.0.0.1:7001"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	cmd := new(proto.Command)
	cmd.Type = proto.Command_COMMAND.Enum()
	cmd.Command = protobuf.String("Hello")
	data, err := protobuf.Marshal(cmd)
	failOnError(err, "Failed to marshal")

	requester, err := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	if err != nil {
		failOnError(err, "Unable to create socket")
	}
	if err = requester.Connect(url); err != nil {
		failOnError(err, "Unable to connect")
	}

	if _, err = requester.SendBytes(data, 0); err != nil {
		failOnError(err, "can't send message on push socket")
	}

	reply, err := requester.Recv(0)
	if err != nil {
		failOnError(err, "Unable to receive reply")
	}
	log.Println(reply)
	
	requester.Close()

	os.Exit(0)
}
