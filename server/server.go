package main

import (
	"log"

	"github.com/emembrives/tinkerings/android-remote/proto"
	protobuf "github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

const (
	zmqURL = "tcp://0.0.0.0:7001"
)

var (
	messages = make(chan []byte)
)

// ServiceProvider exposes endpoints to the ZMQ network.
type ServiceProvider interface {
	Endpoints() []string
	ReadEndpoint(name string) string
	WriteEndpoint(name, value string) string
}

// ZMQServer holds service providers.
type ZMQServer struct {
	Services []ServiceProvider
}

func main() {
	server := ZMQServer{}
	go setupMessageServer(&server)
	websocketListen(&server)
}

func setupMessageServer(server *ZMQServer) {
	responder, err := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	if err != nil {
		panic(err)
	}
	err = responder.Bind(zmqURL)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := responder.RecvBytes(0)
		if err != nil {
			log.Println(err)
			errorMsg := &proto.Response{
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
				ErrorMessage: protobuf.String(err.Error()),
			}
			data, err := protobuf.Marshal(errorMsg)
			logFatalOnError(err)
			responder.SendBytes(data, 0)
			continue
		}

		response := processRequest(incomingRequest, server)
		data, err := protobuf.Marshal(&response)
		logFatalOnError(err)
		responder.SendBytes(data, 0)
	}
}

func processRequest(request *proto.Request, server *ZMQServer) (response proto.Response) {
	switch *request.Type {
	case proto.RequestType_SERVICE_DISCOVERY:
		for _, service := range server.Services {
			for _, endpoint := range service.Endpoints() {
				value := service.ReadEndpoint(endpoint)
				endpointValue := proto.Endpoint{
					Endpoint: protobuf.String(endpoint),
					Value:    protobuf.String(value),
				}
				response.Endpoints = append(response.Endpoints, &endpointValue)
			}
		}
		return response
	case proto.RequestType_WRITE_ENDPOINT:
		// List services
		return
	}
	return
}
