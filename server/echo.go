package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
)

func echo(agent *Agent, p proto.Message) {
	req := p.(*pb.MQEcho)
	rep := &pb.MREcho{}
	rep.Data = proto.String(req.GetData())
	replyMsg(agent, pb.MECHO, rep)
}

func init() {
	registerHandler(pb.MECHO, &pb.MQEcho{}, echo)
}
