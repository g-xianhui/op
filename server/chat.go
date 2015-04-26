package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
)

const (
	_ = iota
	CHAT_PERSONAL
	CHAT_WORLD
	CHAT_GUILD
)

var worldChannel *Broadcast

func chat(agent *Agent, p proto.Message) {
	req := p.(*pb.MQChat)
	chatType := req.GetChatType()
	switch chatType {
	case CHAT_PERSONAL:
		targetId := req.GetTargetId()
		target := agentcenter.find(targetId)
		if target == nil {
			return
		}
		req.From = proto.Uint32(agent.getRoleId())
		sendInnerMsg(target, "redirect", &NetMsgInside{pb.MCHAT, req})
	case CHAT_WORLD:
		req.From = proto.Uint32(agent.getRoleId())
		m := packInnerMsg("redirect", &NetMsgInside{pb.MCHAT, req})
		broadcast(worldChannel, m)
	}
}

func init() {
	worldChannel = createBroadcast()
	registerHandler(pb.MCHAT, &pb.MQChat{}, chat)
}
