package services

import (
	"github.com/emembrives/tinkerings/android-remote/proto"
	protobuf "github.com/golang/protobuf/proto"
)

type PresentationService struct {
}

func NewPresentationService() *PresentationService {
	return new(PresentationService)
}

func (s *PresentationService) GetDescription() *proto.ServiceDefinition {
	p := &proto.ServiceDefinition{}
	p.ServiceName = protobuf.String("presentation")
	p.Method = append(p.Method,
		&proto.MethodDefinition{
			Name: protobuf.String("Next"),
		})
	p.Method = append(p.Method,
		&proto.MethodDefinition{
			Name: protobuf.String("Previous"),
		})
	return p
}

func (s *PresentationService) MakeCall(call *proto.RPCRequest) ([]*proto.ObjectType, error) {
	return nil, nil
}
