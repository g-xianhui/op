package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
	"net"
)

type InnerMsg struct {
	cmd string
	ud  interface{}
}

func (m *InnerMsg) getMsgType() int {
	return MSG_INNER
}

func sendInnerMsg(agent *Agent, cmd string, ud interface{}) {
	m := &InnerMsg{cmd, ud}
	agent.msg <- m
}

type InnerMsgCB func(*Agent, interface{})

func registerInnerMsgHandler(cmd string, cb InnerMsgCB) {
	innerMsgHandlers[cmd] = cb
}

func dispatchInnerMsg(agent *Agent, m *InnerMsg) {
	h, ok := innerMsgHandlers[m.cmd]
	if !ok {
		log(ERROR, "inner msg[%s] handler not found\n", m.cmd)
		return
	}
	h(agent, m.ud)
}

var innerMsgHandlers = map[string]InnerMsgCB{
	"quit":      hAgentQuit,
	"save":      hAgentSave,
	"refresh":   hAgentRefresh,
	"redirect":  hAgentRedirect,
	"worldchat": hAgentWorldChat,
}

func hAgentQuit(agent *Agent, ud interface{}) {
	agent.quit(SERVERCLOSE)
	done := ud.(chan struct{})
	done <- struct{}{}
}

func hAgentSave(agent *Agent, ud interface{}) {
	agent.save()
}

type IMsgRefresh struct {
	conn    net.Conn
	session uint32
}

func hAgentRefresh(agent *Agent, ud interface{}) {
	d := ud.(*IMsgRefresh)
	agent.refresh(d.conn, d.session)
}

type IMsgRedirect struct {
	t uint32
	p proto.Message
}

func hAgentRedirect(agent *Agent, ud interface{}) {
	m := ud.(*IMsgRedirect)
	replyMsg(agent, m.t, m.p)
}

func hAgentWorldChat(agent *Agent, ud interface{}) {
	p := ud.(*pb.MQChat)
	replyMsg(agent, pb.MCHAT, p)
}
