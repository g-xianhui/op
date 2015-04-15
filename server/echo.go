package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
)

const MQECHO = 1

func echo(agent *Agent, p proto.Message) {
	req := p.(*pb.MQEcho)
	log("echo %s", req.GetData())
}

func init() {
	registerHandler(MQECHO, &pb.MQEcho{}, echo)
}
