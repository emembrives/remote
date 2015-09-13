package main

import (
	"log"

	"github.com/emembrives/remote/proto"
	protobuf "github.com/golang/protobuf/proto"

	zmq "github.com/pebbe/zmq4"
)

const (
	zmqURL = "tcp://0.0.0.0:7001"
)

// ServiceProvider exposes endpoints to the ZMQ network.
type ServiceProvider interface {
	Name() string
	Endpoints() []string
	ReadEndpoint(name string) string
	WriteEndpoint(name, value string) string
}

// ZMQServer holds service providers.
type ZMQServer struct {
	services         map[string]ServiceProvider
	AddService       chan ServiceProvider
	RemoveService    chan ServiceProvider
	incomingMessages chan proto.Request
	outgoingMessages chan proto.Response
}

func main() {
	server := &ZMQServer{
		services:         make(map[string]ServiceProvider),
		AddService:       make(chan ServiceProvider),
		RemoveService:    make(chan ServiceProvider),
		incomingMessages: make(chan proto.Request),
		outgoingMessages: make(chan proto.Response),
	}
	go server.setupMessageServer()
	websocketListen(server)
}

func (server *ZMQServer) setupMessageServer() {
	go server.receiveMessages()

	for {
		select {
		case message := <-server.incomingMessages:
			server.outgoingMessages <- server.processRequest(message)
			break
		case service := <-server.AddService:
			server.services[service.Name()] = service
			break
		case service := <-server.RemoveService:
			delete(server.services, service.Name())
			break
		}
	}
}

func (server *ZMQServer) receiveMessages() {
	responder, err := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	logFatalOnError(err)

	err = responder.Bind(zmqURL)
	logFatalOnError(err)

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

		server.incomingMessages <- *incomingRequest
		response := <-server.outgoingMessages
		data, err := protobuf.Marshal(&response)
		logFatalOnError(err)
		responder.SendBytes(data, 0)
	}
}

func (server *ZMQServer) processRequest(request proto.Request) (response proto.Response) {
	switch *request.Type {
	case proto.RequestType_SERVICE_DISCOVERY:
		for _, service := range server.services {
			for _, endpoint := range service.Endpoints() {
				value := service.ReadEndpoint(endpoint)
				endpointValue := proto.Endpoint{
					Service:  protobuf.String(service.Name()),
					Endpoint: protobuf.String(endpoint),
					Value:    protobuf.String(value),
				}
				response.Endpoints = append(response.Endpoints, &endpointValue)
			}
		}
		return response
	case proto.RequestType_WRITE_ENDPOINT:
		for _, endpoint := range request.WriteRequest {
			server.services[*endpoint.Service].WriteEndpoint(*endpoint.Endpoint, *endpoint.Value)
		}
		return
	}
	return
}
