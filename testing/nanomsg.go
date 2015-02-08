package main

import (
	"fmt"
	"log"
	"os"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/sterops/tinkerings/android-remote/proto"

	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/push"
	"github.com/gdamore/mangos/transport/tcp"
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

	var sock mangos.Socket

	if sock, err = push.NewSocket(); err != nil {
		failOnError(err, "can't get new push socket: %s")
	}
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		failOnError(err, "can't dial on push socket")
	}
	if err = sock.Send(data); err != nil {
		failOnError(err, "can't send message on push socket")
	}
	sock.Close()

	os.Exit(0)
}
