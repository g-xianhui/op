package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
)

func hRolelist(agent *Agent, p proto.Message) {
	replyRolelist(agent)
}

func hLogin(agent *Agent, p proto.Message) {
	req := p.(*pb.MQLogin)
	id := req.GetRoleid()
	errno := login(agent, id)
	rep := &pb.MRLogin{}
	rep.Errno = proto.Uint32(errno)
	replyMsg(agent, pb.MRLOGIN, rep)
}

func hCreateRole(agent *Agent, p proto.Message) {
	req := p.(*pb.MQCreateRole)
	rep := &pb.MRCreateRole{}
	basic, errno := createRole(agent, req.GetOcc(), req.GetName())
	rep.Errno = proto.Uint32(errno)
	if basic != nil {
		rep.Basic = toRoleBasic(basic)
	}
	replyMsg(agent, pb.MRCREATEROLE, rep)
}

func init() {
	registerHandler(pb.MQROLELIST, &pb.MQRolelist{}, hRolelist)
	registerHandler(pb.MQLOGIN, &pb.MQLogin{}, hLogin)
	registerHandler(pb.MQCREATEROLE, &pb.MQCreateRole{}, hCreateRole)
}
