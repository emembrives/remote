package main

import (
	"log"

	"github.com/emembrives/tinkerings/android-remote/proto"
	protobuf "github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

const (
	zmqUrl = "tcp://0.0.0.0:7001"
)

var (
	messages = make(chan []byte)
)

func main() {
	go setupMessageServer()
	//websocketListen()
}

func setupMessageServer() {
	responder, err := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	if err != nil {
		panic(err)
	}
	err = responder.Bind(zmqUrl)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := responder.RecvBytes(0)
		if err != nil {
			log.Println(err)
			errorMsg := &proto.Response{
				ErrorType:    proto.Response_UNKNOWN.Enum(),
				ErrorMessage: protobuf.String(err.Error()),
			}
			data, err := protobuf.Marshal(errorMsg)
			logFatalOnError(err)
			responder.SendBytes(data, 0)
			continue
		}

		incomingRequest := new(proto.Request)
		err = protobuf.Unmarshal(msg, incomingRequest)

		if err != nil {
			log.Println(err)
			errorMsg := &proto.Response{
				ErrorType:    proto.Response_UNKNOWN.Enum(),
				ErrorMessage: protobuf.String(err.Error()),
			}
			data, err := protobuf.Marshal(errorMsg)
			logFatalOnError(err)
			responder.SendBytes(data, 0)
			continue
		}

		response := processRequest(incomingRequest)
		data, err := protobuf.Marshal(&response)
		logFatalOnError(err)
		responder.SendBytes(data, 0)
	}
}

func processRequest(request *proto.Request) (response proto.Response) {
	response.Type = request.Type
	response.Host = protobuf.String("luna")
	switch *request.Type {
	case proto.RequestType_PING:
		return
	case proto.RequestType_SERVICES:
		// List services
		return
	}
	return
}
