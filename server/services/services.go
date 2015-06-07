package services

import (
	"github.com/emembrives/tinkerings/android-remote/proto"
)

type Service interface {
	GetDescription() *proto.ServiceDefinition
	MakeCall(call *proto.RPCRequest) []*proto.ObjectType
}
